#!/bin/bash
set -e
#
#
# Gets dist/tarball_dirty created by Goreleaser (all files in root path) and reorganize files in correct path
#
#
PROJECT_PATH=$1

find dist -regex ".*_dirty\.tar.gz" | while read tarball_dirty; do
  echo "tarball_dirty: $tarball_dirty"
  tarball=${tarball_dirty/_dirty.tar.gz} # strip trailing _dirty
  tarball=${tarball/dist\/} # strip leading folder name
  echo "tarball: $tarball"  
  TARBALL_CLEAN="${tarball}.tar.gz"
  TARBALL_TMP="dist/tarball_temp"
  echo "TARBALL_CLEAN: $TARBALL_CLEAN"
  TARBALL_CONTENT_PATH="${TARBALL_TMP}/${tarball}_content"
  mkdir -p ${TARBALL_CONTENT_PATH}/var/db/newrelic-infra/newrelic-integrations/bin/
  mkdir -p ${TARBALL_CONTENT_PATH}/etc/newrelic-infra/integrations.d/
  mkdir -p ${TARBALL_CONTENT_PATH}/etc/newrelic-infra/logging.d/
  echo "===> Decompress ${tarball} in ${TARBALL_CONTENT_PATH}"
  tar -xvf ${tarball_dirty} -C ${TARBALL_CONTENT_PATH}

  echo "===> Move files inside ${tarball}"
  mv ${TARBALL_CONTENT_PATH}/nri-${INTEGRATION} "${TARBALL_CONTENT_PATH}/var/db/newrelic-infra/newrelic-integrations/bin/"
  mv ${TARBALL_CONTENT_PATH}/${INTEGRATION}-definition.yml ${TARBALL_CONTENT_PATH}/var/db/newrelic-infra/newrelic-integrations/
  mv ${TARBALL_CONTENT_PATH}/${INTEGRATION}-config.yml.sample ${TARBALL_CONTENT_PATH}/etc/newrelic-infra/integrations.d/
  mv ${TARBALL_CONTENT_PATH}/${INTEGRATION}-log.yml.example ${TARBALL_CONTENT_PATH}/etc/newrelic-infra/logging.d/

  echo "===> Creating tarball ${TARBALL_CLEAN}"
  cd ${TARBALL_CONTENT_PATH}
  tar -czvf ../${TARBALL_CLEAN} .
  cd $PROJECT_PATH
  echo "===> Moving tarball ${TARBALL_CLEAN}"
  mv "${TARBALL_TMP}/${TARBALL_CLEAN}" dist/
  echo "===> Cleaning dirty tarball ${tarball_dirty}"
  rm ${tarball_dirty}
done
