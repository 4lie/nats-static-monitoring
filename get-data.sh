#!/usr/bin/env bash

curl http://localhost:9200/core_varz/_search?pretty=true&q=*:*
