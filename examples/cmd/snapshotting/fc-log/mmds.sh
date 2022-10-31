
# check mmds
curl -s --unix-socket /tmp/FCGoSDKSnapshotExample1998987509/fc-0.create http://localhost/mmds

# stop
curl --unix-socket /tmp/FCGoSDKSnapshotExample3006451316/fc-0.create -i \
  -X PUT 'http://localhost/actions'       \
  -H  'Accept: application/json'          \
  -H  'Content-Type: application/json'    \
  -d '{
      "action_type": "SendCtrlAltDel"
   }'


/tmp/FCGoSDKSnapshotExample1495375699/fc-0

curl --unix-socket /tmp/FCGoSDKSnapshotExample1998987509/fc-0.create -i \
    -X PUT "http://localhost/mmds"            \
    -H "Content-Type: application/json"       \
    -d '{
            "latest": {
                  "meta-data": {
                       "ami-id": "ami-12345678",
                       "reservation-id": "r-fea54097"
                  }
            }
    }'

curl -s -H "Accept: application/json" "http://169.254.169.254/latest"

ds=nocloud-net;s=http://169.254.169.254/latest/ network-config=dmVyc2lvbjogMgpldGhlcm5ldHM6CiAgZXRoMDoKICAgIGFkZHJlc3NlczoKICAgIC0gMTY5LjI1NC4xNjkuMjU0LzE2CiAgICBtYXRjaDoKICAgICAgbWFjYWRkcmVzczogQUE6RkM6YWM6MWE6MDA6MDIKICBldGgxOgogICAgYWRkcmVzc2VzOgogICAgLSAxNzIuMjYuMC4zLzI0CiAgICBnYXRld2F5NDogMTcyLjI2LjAuMgogICAgbWF0Y2g6CiAgICAgIG1hY2FkZHJlc3M6IEFBOkZDOmFjOjFhOjAwOjAzCiAgICBuYW1lc2VydmVyczoKICAgICAgYWRkcmVzc2VzOgogICAgICAtIDEuMS4xLjEKICAgICAgLSA4LjguOC44Cg==

https://askubuntu.com/questions/1104285/how-do-i-reload-network-configuration-with-cloud-init

https://askubuntu.com/questions/1405294/ubuntu-20-04-cloud-init-wont-configure-network
