#!/bin/bash

# Ensure that environment variable is set
if [[ -z "${ENDPOINT_URL}" || -z "${START_BLOCK}" || -z "${CHAIN_ID}" ]]; then
  echo "Error: ENDPOINT_URL, START_BLOCK, CHAIN_ID environment variable must be set."
  exit 1
fi

# Specify the file name to be modified
FILE_NAME="project.yaml"  # Change this to your actual file name

# Using sed to replace the target value with the new value in the specified file
sed -i "s#ENDPOINT_URL_ENV_VAR_REF#${ENDPOINT_URL}#g" "${FILE_NAME}"
sed -i "s#CHAIN_ID_ENV_VAR_REF#'${CHAIN_ID}'#g" "${FILE_NAME}"
sed -i "s#-1#${START_BLOCK}#g" "${FILE_NAME}"

if [ $? -eq 0 ]; then
  echo "Endpoint url replacement successful."
else
  echo "An error occurred during the replacement."
  exit 1
fi
