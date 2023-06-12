scout
=====

Give scout a repository, it'll figure out which directories are Jenkins plugins and recommend a java version


building
--------

    go build

usage
-----

    $ ./scout --directory path/to/repo --outfile path/to/desired/result.json
    $ cat path/to/desired/result.json
    {
        "modules":[{
            "path":"",
            "pluginParentVersion":"4.40",
            "pluginBoms":[{
                "name":"bom-2.303.x",
                "version":"1342.v729ca_3818e88"
            }],
            "jenkinsVersion":"2.303.3",
            "recommendedJava":"8"
        }]
    }


Implemented in golang as a learning exercise (feedback welcome!)
