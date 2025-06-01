set -x

curl -X POST localhost:8080/box -d '{"content": "test content @daily_cnt", "author": "author?? probably me"}'
