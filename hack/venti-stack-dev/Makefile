template:
	helm template -n venti-stack-dev vs .

image:
	@helm template -n venti-stack-dev vs . | grep -oP 'image: "?\K[^"]*' | sort | uniq

apply:
	helm upgrade --install --create-namespace -n venti-stack-dev vs .

delete:
	helm uninstall -n venti-stack-dev vs

dep:
	helm dependency update
	helm dependency build
