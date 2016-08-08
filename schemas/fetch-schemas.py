#!/usr/bin/env python
from lxml import html
import requests
import json

page = requests.get('https://docs.atlassian.com/jira/REST/cloud')
tree = html.fromstring(page.content)

schemas = tree.xpath("//div[@class='representation-doc-block']//code/text()")

for schema in schemas:
    try:
        data = json.loads(schema)
        if "title" in data:
            title = data["title"].replace(" ", "")
            print "Writing {}.json".format(title)
            with open("{}.json".format(title), 'w') as f:
                f.write(schema)
    except:
        True

