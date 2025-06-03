for i in $(seq 1 10); do
    curl 127.0.0.1:8080/like/2 -X POST
done
