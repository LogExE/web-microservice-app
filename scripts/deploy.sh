podman kube play pod.yaml --publish 8080:8080
podman exec web-go-broker /opt/kafka/bin/kafka-topics.sh --bootstrap-server localhost:9092 --create --topic newBoxes
podman exec web-go-broker /opt/kafka/bin/kafka-topics.sh --bootstrap-server localhost:9092 --create --topic likeBoxes
