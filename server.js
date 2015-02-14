require( 'parsoid/lib/core-upgrade.js' );

var port = 3000;

var cluster = require('cluster');
var numCPUs = require('os').cpus().length;

if (cluster.isMaster) {
    // Fork workers
    for (var i = 0; i < numCPUs; i++) {
        cluster.fork();
    }

    cluster.on('exit', function(worker, code, signal) {
        console.log('worker ' + worker.process.pid + ' died');
    });
} else {
    var express = require('express');
    var bodyParser = require('body-parser');
    var app = express();
    app.use(bodyParser.text({limit: '50mb'}));

    // wikitext to HTML parsing using Parsoid: code taken from parser.js test program included
    // in Parsoid project's source.
    var ParserEnv = require('parsoid/lib/mediawiki.parser.environment.js').MWParserEnvironment;
    var ParsoidConfig = require('parsoid/lib/mediawiki.ParsoidConfig.js').ParsoidConfig;
    var DU = require('parsoid/lib/mediawiki.DOMUtils.js').DOMUtils;

    var getParserEnv = Promise.promisify(ParserEnv.getParserEnv, false, ParserEnv);

    // test with:
    // curl -X POST -v -H"Content-Type:text/plain" \
    //      --data “stuff that goes in POST body" \
    //      localhost:3000/parse
    //
    var prefix = "enwiki";
    app.post('/parse', function(req, res) {
        // Convert wikitext to HTML and send it as response
        var parsoidConfig = new ParsoidConfig(null, null);
        //parsoidConfig.fetchTemplates = false;
        getParserEnv(parsoidConfig, null, prefix, null, null).then(function(env) {
            return new Promise(function(resolve) {
                var parser = env.pipelineFactory.getPipeline('text/x-mediawiki/full');
                parser.once('document', resolve);
                // Kick off the pipeline by feeding the input into the parser pipeline
                env.setPageSrcInfo(req.body);
                parser.processToplevelDoc(env.page.src);
            }).then(function(doc) {
                res.send(DU.serializeNode(doc));
            }).done();
        });
    });

    app.listen(port);
    console.log('Listening on port ' + port);
}
