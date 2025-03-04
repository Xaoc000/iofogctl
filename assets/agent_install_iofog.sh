#!/bin/sh
set -x
set -e

do_stop_iofog() {
	if command_exists iofog-agent; then
		sudo service iofog-agent stop
	fi
}

do_check_iofog_on_arm() {
    if [ "$lsb_dist" = "raspbian" ] || [ "$(uname -m)" = "armv7l" ] || [ "$(uname -m)" = "aarch64" ] || [ "$(uname -m)" = "armv8" ]; then

        echo "# We re on ARM ($(uname -m)) : Updating config.xml to use correct docker_url"
		$sh_c 'sed -i -e "s|<docker_url>.*</docker_url>|<docker_url>tcp://127.0.0.1:2375/</docker_url>|g" /etc/iofog-agent/config.xml'

		echo "# Restarting iofog-agent service"
		$sh_c "service iofog-agent stop"
		sleep 3
		$sh_c "service iofog-agent start"
	fi
}

do_install_iofog() {
	AGENT_CONFIG_FOLDER=/etc/iofog-agent
	SAVED_AGENT_CONFIG_FOLDER=/tmp/agent-config-save
	echo "# Installing ioFog agent..."

	# Save iofog-agent config
	if [ -d ${AGENT_CONFIG_FOLDER} ]; then
		sudo rm -rf ${SAVED_AGENT_CONFIG_FOLDER}
		sudo mkdir -p ${SAVED_AGENT_CONFIG_FOLDER}
		sudo cp -r ${AGENT_CONFIG_FOLDER}/* ${SAVED_AGENT_CONFIG_FOLDER}/
	fi

	prefix=$([ -z "$token" ] && echo "" || echo "$token:@")

	case "$lsb_dist" in
		ubuntu)
			curl -s "https://${prefix}packagecloud.io/install/repositories/$repo/script.deb.sh" | $sh_c "bash"
			$sh_c "apt-get install -y --allow-downgrades iofog-agent=$agent_version"
			;;
		fedora|centos)
			curl -s "https://${prefix}packagecloud.io/install/repositories/$repo/script.rpm.sh" | $sh_c "bash"
			$sh_c "yum install -y iofog-agent-"$agent_version"-1.noarch"
			;;
		debian|raspbian)
			curl -s "https://${prefix}packagecloud.io/install/repositories/$repo/script.deb.sh" | $sh_c "bash"
			$sh_c "apt-get install -y --allow-downgrades iofog-agent=$agent_version"
			;;
	esac

	do_check_iofog_on_arm

	# Restore iofog-agent config
	if [ -d ${SAVED_AGENT_CONFIG_FOLDER} ]; then
		sudo mv ${SAVED_AGENT_CONFIG_FOLDER}/* ${AGENT_CONFIG_FOLDER}/
		sudo rmdir ${SAVED_AGENT_CONFIG_FOLDER}
	fi
}

agent_version="$1"
repo=$([ -z "$2" ] && echo "iofog/iofog-agent" || echo "$2")
token="$3"
echo "Using variables"
echo "version: $agent_version"
echo "repo: $repo"
echo "token: $token"

. /tmp/agent_init.sh
init
do_stop_iofog
do_install_iofog