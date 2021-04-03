dev:
	mkdir -p ${PWD}/bin/tmp
	go build -o ${PWD}/bin/tmp/terraform-provider-windns
	chmod +x ${PWD}/bin/tmp/terraform-provider-windns
	echo "provider_installation {\n  dev_overrides {\n    \"localhost/bhoriuchi/windns\" = \"${PWD}/bin/tmp\"\n  }\n  direct {}\n}" > ${PWD}/dev.tfrc
	echo "export TF_CLI_CONFIG_FILE=${PWD}/dev.tfrc" > ${PWD}/dev.env