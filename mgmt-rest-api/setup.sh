#!/bin/bash

e(){
    echo "$(date): $*"
}

NEEDED_BASE_PYTHON=3.10.3
NEEDED_PYENV_ENV=dns_mgmt_rest_api

if ! which pyenv >/dev/null 2>&1; then
    e "ERROR: must have pyenv installed"
    exit 1
fi

if ! pyenv versions |grep -q "${NEEDED_BASE_PYTHON}"; then
    e "Building python ${NEEDED_BASE_PYTHON}"
    e pyenv install $NEEDED_BASE_PYTHON
fi

if ! pyenv local | grep -q $NEEDED_PYENV_ENV; then
    e "Installing needed ${NEEDED_PYENV_ENV}"
    pyenv virtualenv $NEEDED_BASE_PYTHON $NEEDED_PYENV_ENV
    e "Adding ${NEEDED_PYENV_ENV} as local"
    pyenv local $NEEDED_PYENV_ENV

fi

if ! pip list | grep -q gql; then
    pip install pip --upgrade
    e "Adding needed pip modules"
    pip install -r requirements.txt
fi
