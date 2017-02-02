package googlecompute

const StartupScriptKey string = "startup-script"
const StartupScriptStatusKey string = "startup-script-status"
const StartupWrappedScriptKey string = "packer-wrapped-startup-script"

const StartupScriptStatusDone string = "done"
const StartupScriptStatusError string = "error"
const StartupScriptStatusNotDone string = "notdone"

var StartupScriptLinux string = `#!/usr/bin/env bash
echo "Packer startup script starting."

RETVAL=0
BASEMETADATAURL=http://metadata/computeMetadata/v1/instance/
GCEPROFILE=/etc/profile.d/google-cloud-sdk.sh

if [ -f $GCEPROFILE ]; then
    shopt -s expand_aliases
    source $GCEPROFILE
fi

GetMetadata () {
  echo "$(curl -f -H "Metadata-Flavor: Google" ${BASEMETADATAURL}/${1} 2> /dev/null)"
}

ZONE=$(GetMetadata zone | grep -oP "[^/]*$")
NAME=$(GetMetadata name)

SetMetadata () {
  gcloud compute instances add-metadata ${NAME} --metadata ${1}=${2} --zone ${ZONE}
}

STARTUPSCRIPT=$(GetMetadata attributes/packer-wrapped-startup-script)
STARTUPSCRIPTPATH=$(mktemp -d)/packer-wrapped-startup-script
if [ -f "/var/log/startupscript.log" ]; then
  STARTUPSCRIPTLOGPATH=/var/log/startupscript.log
else
  STARTUPSCRIPTLOGPATH=/var/log/daemon.log
fi
STARTUPSCRIPTLOGDEST=$(GetMetadata attributes/startup-script-log-dest)

if [[ ! -z $STARTUPSCRIPT ]]; then
  echo "Executing user-provided startup script..."
  echo "${STARTUPSCRIPT}" > ${STARTUPSCRIPTPATH}
  chmod +x ${STARTUPSCRIPTPATH}
  ${STARTUPSCRIPTPATH}
  RETVAL=$?

  if [[ ! -z $STARTUPSCRIPTLOGDEST ]]; then
    echo "Uploading user-provided startup script log to ${STARTUPSCRIPTLOGDEST}..."
    gsutil -h "Content-Type:text/plain" cp ${STARTUPSCRIPTLOGPATH} ${STARTUPSCRIPTLOGDEST}
  fi

  rm ${STARTUPSCRIPTPATH}
fi

echo "Packer startup script done."
SetMetadata startup-script-status done
exit $RETVAL`

var StartupScriptWindows string = ""
