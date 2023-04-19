```json
GET _search
{
  "query": {
    "match_all": {}
  }
}

#获取所有索引
GET _cat/indices

#多次操作会生成多个数据
POST user/_doc
{
  "name":"bobby",
  "company":"imooc"
}

POST user/_create/1
{
  "name":"bobby",
  "company":"imooc"
}

#添加数据 使用PUT 则ID必须存在
PUT account/_doc/1
{
  "name":"bobby",
  "age":18,
  "company":[
    {
      "name":"imooc",
      "address":"beijing"
    },
    {
      "name":"imooc2",
      "address":"shanghai"
    }
  ]
}

GET account

#获取数据
GET user/_doc/1

#获取具体数据
GET user/_source/1

GET _search?q=bobby

GET user/_search
{
  "query": {
    "match_all": {}
  }
}

GET account/_search
{
  "query": {"match_all": {}}
}

POST user/_doc/2
{
  "name":"bobby2",
  "company":"imooc"
}
#会把之前的所有数据都删除 -- 全部替换
#即使数据是一样的也会进行替换
POST user/_doc/2
{
  "age":18
}

PUT user/_doc/2
{
  "name":"bobby"
}

#这是条件更新 不会删除之前数据
POST user/_update/2
{
  "doc":{
    "age":18
  }
}

GET user/_doc/2

#删除数据
DELETE user/_doc/2

#删除索引
DELETE user

#批量操作
POST /_bulk
{"index":{"_index":"user","_id":"3"}}
{"name":"bobby"}
{"delete":{"_index":"user","_id":"2"}}

GET bank/_search
{
  "query": {"match_all": {}}
}

#组合结果数据
GET /_mget
{
  "docs":[
    {
      "_index":"bank",
      "_id":"1"
    },
    {
      "_index":"account",
      "_id":"1"
    }
  ]
}

#对于es来说 from和size分页在数据量比较小的情况下可行
GET bank/_search
{
  "query": {
    "match_all": {}
  },
  "from":4,
  "size": 4
}

#大小写不敏感 默认大写全部转换成小写
GET bank/_search
{
  "query": {
    "match": {
      "address": "Madison street"
    }
  }
}

#短语查询
GET bank/_search
{
  "query": {
    "match_phrase": {
      "address": "Madison street"
    }
  }
}

#多字段查询
GET bank/_search
{
  "query": {
    "multi_match": {
      "query": "street",
      "fields": ["address","firstname"]
    }
  }
}

#query_string查询 Bristol AND Street AND表示连接符
GET bank/_search
{
  "query": {
    "query_string": {
      "default_field": "address", 
      "query": "Madison AND street"
    }
  }
}

GET bank/_search
{
  "query": {
    "match_all": {}
  }
}

GET bank/_search
{
  "query": {
    "term": {
      "address": {
        "value": "street"
      }
    }
  }
}

#范围搜索
GET bank/_search
{
  "query": {
    "range": {
      "age": {
        "gte": 10,
        "lte": 20
      }
    }
  }
}

POST bank/_doc
{
  "name":"bobby"
}

#查询所有有age字段的数据
GET bank/_search
{
  "query": {
    "exists": {
      "field": "age"
    }
  }
}

POST bank/_doc
{
  "school":"middle school"
}

#查询所有有school字段的数据
GET bank/_search
{
  "query": {
    "exists": {
      "field": "school"
    }
  }
}


#模糊查询
GET bank/_search
{
  "query": {
    "fuzzy": {
      "address": "street"
    }
  }
}

#大小写不敏感
GET bank/_search
{
  "query": {
    "match": {
      "address": {
        "query": "Bristol Street",
        "fuzziness": 1
      }
    }
  }
}

# bool 组合 mush should mush_not filter
GET bank/_search
{
  "query": {
    "bool": {
      "must": [
        {
          "term": {
            "state": {
              "value": "tn"
            }
          }
        },
        {
          "range": {
            "age": {
              "gte": 20,
              "lte": 30
            }
          }
        }
      ],
      "must_not": [
        {
          "term": {
            "gender": {
              "value": "m"
            }
          }
        }
      ],
      "should": [
        {
          "term": {
            "firstname": {
              "value": "decker"
            }
          }
        }
      ],
      "filter": [
        {
          "range": {
            "age": {
              "gte": 25,
              "lte": 30
            }
          }
        }
      ]
    }
  }
}

GET bank/_mapping

#使用keyword类型进行精确匹配
GET bank/_search
{
  "query": {
    "match": {
      "address.keyword": "671 Bristol Street"
    }
  }
}

PUT usertest
{
  "mappings": {
    "properties": {
      "age":{
        "type": "integer"
      },
      "name":{
        "type": "text"
      },
      "desc":{
        "type": "keyword"
      }
    }
  }
}

GET usertest/_mapping

POST usertest/_doc
{
  "age":18,
  "name":"Hattie Bond",
  "desc":"671 Bristol Street"
}

GET usertest/_search
{
  "query": {
    "match": {
      "name": "Bond"
    }
  }
}

GET usertest/_search
{
  "query": {
    "match": {
      "desc": "671 Bristol Street"
    }
  }
}

POST cn/_doc
{
  "name":"中华牙膏"
}

GET cn/_search
{
  "query": {
    "match": {
      "name": "中立"
    }
  }
}

GET _analyze
{
  "text":"中国科学技术大学",
  "analyzer": "ik_max_word"
}

DELETE cn

PUT cn
{
  "mappings": {
    "properties": {
      "name":{
        "type": "text",
        "analyzer": "ik_max_word"
      }
    }
  }
}

GET _analyze
{
  "text":"中国科学技术大学",
  "analyzer": "ik_smart"
}


POST cn/_doc
{
  "name":"中国科学技术大学"
}

GET cn/_search
{
  "query": {
    "match": {
      "name": "大学"
    }
  }
}

PUT newCn
{
  "mappings": {
    "properties": {
      "name":{
        "type": "text",
        "analyzer": "ik_smart",
        "search_analyzer": "ik_smart"
      }
    }
  }
}

GET cn

GET _analyze
{
  "analyzer": "ik_smart",
  "text": ["慕课网的课程资源非常丰富"]
}

GET mybank/_search
{
  "query": {"match_all": {}}
}


GET mygoods/_mapping

GET _cat/indices

GET goods


```