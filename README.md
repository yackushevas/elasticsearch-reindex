# ElasticSearch-Reindex

ElasticSearch-Reindex is an simple utility for reindex matched pattern index.

[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/yackushevas/elasticsearch-reindex/master/LICENSE)

**Requirements**

Elasticsearch 5.X

**Example:**

reindex one index
```sh
$ ./elasticsearch-reindex -u http://localhost:9200 -i logstash.nginx-2017.02.25
```

reindex by pattern
```sh
$ ./elasticsearch-reindex -u http://localhost:9200 -i logstash.nginx*
```

**Algorithm**

1. Setting index logstash.nginx-2017.02.25 to read-only
2. Reindex index logstash.nginx-2017.02.25
3. Wait for completion task
4. Deleting index logstash.nginx-2017.02.25
5. Adding logstash.nginx-2017.02.25 alias to index logstash.nginx-2017.02.25-20170202550405 (add timestamp)
