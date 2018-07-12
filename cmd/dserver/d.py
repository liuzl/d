import json
with open('input.txt') as fd:
    for line in fd:
        line = line.strip()
        item = line.split('\t')
        rec = {'k':item[0], 'v':{item[1]:item[2]}}
        print(json.dumps(rec, ensure_ascii=False))
