#!/bin/bash

# template code from https://github.com/msound/localmongo/tree/master
echo "sleeping for 10 seconds"
sleep 10

# For MacOS, add the following to /etc/hosts
# 127.0.0.1 mongo1
# 127.0.0.1 mongo2
# 127.0.0.1 mongo3

echo mongo_setup.sh time now: `date +"%T" `
mongo --host mongo1:27017 <<EOF
  var cfg = {
    "_id": "rs",
    "version": 1,
    "members": [
      {
        "_id": 0,
        "host": "mongo1:27017",
        "priority": 2
      },
      {
        "_id": 1,
        "host": "mongo2:27017",
        "priority": 0
      },
      {
        "_id": 2,
        "host": "mongo3:27017",
        "priority": 0
      }
    ]
  };
  rs.initiate(cfg);
  rs.slaveOk();
  db.getMongo().setReadPref('nearest');
  db.getMongo().setSlaveOk();
EOF
