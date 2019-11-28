// Copyright © 2019 The Things Network Foundation, The Things Industries B.V.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import React from 'react'

import SubmitButton from '../../../../components/submit-button'
import SubmitBar from '../../../../components/submit-bar'
import Input from '../../../../components/input'
import Checkbox from '../../../../components/checkbox'
import Form from '../../../../components/form'

import diff from '../../../../lib/diff'
import m from '../../../components/device-data-form/messages'
import PropTypes from '../../../../lib/prop-types'
import sharedMessages from '../../../../lib/shared-messages'

import { parseLorawanMacVersion, hasExternalJs } from '../utils'
import validationSchema from './validation-schema'

// The Join Server can store end device fields while not exposing the root keys. This means
// that the `root_keys` object is present while `root_keys.nwk_key` == nil or `root_keys.app_key == nil`
// must hold. See https://github.com/TheThingsNetwork/lorawan-stack/issues/1473
const isNwkKeyHidden = ({ root_keys }) => Boolean(root_keys) && !Boolean(root_keys.nwk_key)
const isAppKeyHidden = ({ root_keys }) => Boolean(root_keys) && !Boolean(root_keys.app_key)

const JoinServerForm = React.memo(props => {
  const { device, onSubmit } = props

  const isNewLorawanVersion = parseLorawanMacVersion(device.lorawan_version) >= 110
  const externalJs = hasExternalJs(device)

  const formRef = React.useRef(null)
  const [error, setError] = React.useState('')
  const [resetsJoinNonces, setResetsJoinNonces] = React.useState(device.resets_join_nonces)

  // Setup and memoize initial form state.
  const initialValues = React.useMemo(() => {
    const extJs = hasExternalJs(device)
    const {
      root_keys = {
        nwk_key: {},
        app_key: {},
      },
      resets_join_nonces,
      lorawan_version,
      net_id = '',
    } = device

    return {
      resets_join_nonces,
      root_keys,
      _external_js: hasExternalJs(device),
      _lorawan_version: lorawan_version,
      net_id: extJs ? undefined : net_id,
    }
  }, [device])

  // Setup and memoize callbacks for changes to `resets_join_nonces` for displaying the field warning.
  const handleResetsJoinNoncesChange = React.useCallback(
    evt => {
      setResetsJoinNonces(evt.target.checked)
    },
    [setResetsJoinNonces],
  )

  const onFormSubmit = React.useCallback(
    async (values, { setSubmitting, resetForm }) => {
      const castedValues = validationSchema.cast(values)
      const updatedValues = diff(initialValues, castedValues, ['_external_js', '_lorawan_version'])

      setError('')
      try {
        await onSubmit(updatedValues)
        resetForm(castedValues)
      } catch (err) {
        setSubmitting(false)
        setError(err)
      }
    },
    [initialValues, onSubmit],
  )

  const nwkKeyHidden = isNwkKeyHidden(device)
  const appKeyHidden = isAppKeyHidden(device)

  let appKeyPlaceholder = m.leaveBlankPlaceholder
  if (externalJs) {
    appKeyPlaceholder = sharedMessages.provisionedOnExternalJoinServer
  } else if (appKeyHidden) {
    appKeyPlaceholder = m.unexposed
  }

  let nwkKeyPlaceholder = m.leaveBlankPlaceholder
  if (externalJs) {
    nwkKeyPlaceholder = sharedMessages.provisionedOnExternalJoinServer
  } else if (nwkKeyHidden) {
    nwkKeyPlaceholder = m.unexposed
  }

  return (
    <Form
      validationSchema={validationSchema}
      initialValues={initialValues}
      onSubmit={onFormSubmit}
      formikRef={formRef}
      error={error}
      enableReinitialize
    >
      <Form.Field
        title={m.netID}
        description={m.netIDDescription}
        name="net_id"
        type="byte"
        min={3}
        max={3}
        component={Input}
        disabled={externalJs}
      />
      <Form.Field
        title={sharedMessages.appKey}
        name="root_keys.app_key.key"
        type="byte"
        min={16}
        max={16}
        placeholder={appKeyPlaceholder}
        description={m.appKeyDescription}
        component={Input}
        disabled={externalJs || appKeyHidden}
      />
      {isNewLorawanVersion && (
        <Form.Field
          title={sharedMessages.nwkKey}
          name="root_keys.nwk_key.key"
          type="byte"
          min={16}
          max={16}
          placeholder={nwkKeyPlaceholder}
          description={m.nwkKeyDescription}
          component={Input}
          disabled={externalJs || nwkKeyHidden}
        />
      )}
      {isNewLorawanVersion && (
        <Form.Field
          title={m.resetsJoinNonces}
          onChange={handleResetsJoinNoncesChange}
          warning={resetsJoinNonces ? m.resetWarning : undefined}
          name="resets_join_nonces"
          component={Checkbox}
          disabled={externalJs}
        />
      )}
      <Form.Field
        title={sharedMessages.macVersion}
        name="_lorawan_version"
        component={Input}
        type="hidden"
        hidden
        disabled
      />
      <SubmitBar>
        <Form.Submit component={SubmitButton} message={sharedMessages.saveChanges} />
      </SubmitBar>
    </Form>
  )
})

JoinServerForm.propTypes = {
  device: PropTypes.device.isRequired,
  onSubmit: PropTypes.func.isRequired,
}

export default JoinServerForm
