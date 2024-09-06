#!/bin/bash

GEN_PKG_NAME=apigen
API_GEN_DIR=./internal/apigen
API_MODEL_DIR=${API_GEN_DIR}/apimodel
API_INTERFACES_DIR=${API_GEN_DIR}/interfaces

# Generate server implementation for OpenAPI
openapi-generator-cli generate \
		-i ./limepipes-api/openapi.yaml \
		-g go-gin-server \
		-o ${API_GEN_DIR} \
		--additional-properties=packageName=${GEN_PKG_NAME},interfaceOnly=true

# create the generated model files to directory and
# change to the correct package name
mkdir -p ${API_MODEL_DIR}
mv ${API_GEN_DIR}/go/model_* ${API_MODEL_DIR}

for model_file in "${API_MODEL_DIR}"/*.go; do
  sed -i "s/${GEN_PKG_NAME}/apimodel/g" "$model_file"
done

# Move the generated server interface file and rename it accordingly
mkdir -p ${API_INTERFACES_DIR}
API_HANDLER_FILE=${API_INTERFACES_DIR}/api_handler.go
mv ${API_GEN_DIR}/go/api_default.go ${API_HANDLER_FILE}
sed -i "s/${GEN_PKG_NAME}/interfaces/g" "${API_HANDLER_FILE}"
sed -i "s/DefaultAPI/ApiHandler/g" "${API_HANDLER_FILE}"

# Move the routers file and rename used API handler variables
API_ROUTER_FILE=${API_GEN_DIR}/routers.go
mv ${API_GEN_DIR}/go/routers.go ${API_ROUTER_FILE}
sed -i "s/DefaultAPI/ApiHandler/g" "${API_ROUTER_FILE}"
sed -i "s/ApiHandler ApiHandler/ApiHandler interfaces.ApiHandler/g" "${API_ROUTER_FILE}"
sed -i '/import (/a\\t"'"github.com/tomvodi/limepipes/internal/${GEN_PKG_NAME}/interfaces"'"' "${API_ROUTER_FILE}"

# Modify data types in model

# Change types to pointer where necessary
sed -i 's/Tunes \[\]ImportTune/Tunes \[\]\*ImportTune/g' ${API_MODEL_DIR}/model_import_file.go
sed -i 's/Set BasicMusicSet/Set \*BasicMusicSet/g' ${API_MODEL_DIR}/model_import_tune.go

# Change string Ids to uuid.UUID
for model_file in "${API_MODEL_DIR}"/*.go; do
  sed -i "s/Id string/Id uuid.UUID/g" "$model_file"
  sed -i "s/Tunes \[\]string/Tunes \[\]uuid.UUID/g" "$model_file"

  # add uuid import if necessary
  if grep -q uuid\.UUID "$model_file"; then
    sed -i '/package apimodel/a import "github.com/google/uuid"' "$model_file"
  fi
done

# Add necessary binding struct tags
# Update models
gomodifytags -w -file ${API_MODEL_DIR}/model_update_set.go -struct UpdateSet -field Title \
  -add-tags binding -template "required"
gomodifytags -w -file ${API_MODEL_DIR}/model_update_tune.go -struct UpdateTune -field Title \
  -add-tags binding -template "required"

# Create models
gomodifytags -w -file ${API_MODEL_DIR}/model_create_set.go -struct CreateSet -field Title \
  -add-tags binding -template "required"
gomodifytags -w -file ${API_MODEL_DIR}/model_create_tune.go -struct CreateTune -field Title \
  -add-tags binding -template "required"

# Tune models
gomodifytags -w -file ${API_MODEL_DIR}/model_tune.go -struct Tune -field Id \
  -add-tags binding -template "required"
gomodifytags -w -file ${API_MODEL_DIR}/model_tune.go -struct Tune -field Title \
  -add-tags binding -template "required"
gomodifytags -w -file ${API_MODEL_DIR}/model_import_tune.go -struct ImportTune -field Id \
  -add-tags binding -template "required"
gomodifytags -w -file ${API_MODEL_DIR}/model_import_tune.go -struct ImportTune -field Title \
  -add-tags binding -template "required"

# Set models
gomodifytags -w -file ${API_MODEL_DIR}/model_basic_music_set.go -struct BasicMusicSet -field Id \
  -add-tags binding -template "required"
gomodifytags -w -file ${API_MODEL_DIR}/model_basic_music_set.go -struct BasicMusicSet -field Title \
  -add-tags binding -template "required"
gomodifytags -w -file ${API_MODEL_DIR}/model_music_set.go -struct MusicSet -field Id \
  -add-tags binding -template "required"
gomodifytags -w -file ${API_MODEL_DIR}/model_music_set.go -struct MusicSet -field Title \
  -add-tags binding -template "required"

gomodifytags -w -file ${API_MODEL_DIR}/model_import_file.go -struct ImportFile -field Name \
  -add-tags binding -template "required"

# Remove unnecessary files and directories
rm  ${API_GEN_DIR}/go.mod  ${API_GEN_DIR}/main.go \
    ${API_GEN_DIR}/.openapi-generator-ignore \
    ${API_GEN_DIR}/Dockerfile
rm -rf ${API_GEN_DIR}/api
rm -rf ${API_GEN_DIR}/.openapi-generator
rm -rf ${API_GEN_DIR}/go
rm -f ${API_MODEL_DIR}/model_basic_set_properties.go
rm -f ${API_MODEL_DIR}/model_create_update_set_properties.go
rm -f ${API_MODEL_DIR}/model_health_200_*.go
rm -f ${API_MODEL_DIR}/model_health_status.go
rm -f ${API_MODEL_DIR}/model_import_info.go
rm -f ${API_MODEL_DIR}/model_object_id.go
rm -f ${API_MODEL_DIR}/model_set_properties.go
rm -f ${API_MODEL_DIR}/model_tune_properties.go
