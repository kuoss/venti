template:
	helm template -n venti-dev vd .

image:
	@helm template -n venti-dev vd . | grep -oP 'image: "?\K[^"]*' | sort | uniq

apply:
	helm upgrade --install --create-namespace -n venti-dev vd .

delete:
	helm uninstall -n venti-dev vd
