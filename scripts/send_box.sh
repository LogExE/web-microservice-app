
echo "Sending example boxes to REST service (and showing responses):"
curl -X POST localhost:8080/box -d '{"content": "test content 123", "author": "author?? probably me"}'
echo
curl -X POST localhost:8080/box -d '{"content": "Haste make waste.", "author": "Konq Giu"}'
echo
curl -X POST localhost:8080/box -d '{"content": "https://some-spammy-site.com", "author": "the hax0r"}'
echo

echo "Boxes we have:"
curl localhost:8080/boxes | jq .
