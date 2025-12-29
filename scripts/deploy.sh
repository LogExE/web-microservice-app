
podman kube play --replace pod.yaml

echo "Creating Kafka topics"
podman exec web-go-broker /opt/kafka/bin/kafka-topics.sh --bootstrap-server localhost:9092 --create --topic newBoxes
podman exec web-go-broker /opt/kafka/bin/kafka-topics.sh --bootstrap-server localhost:9092 --create --topic likeBoxes
