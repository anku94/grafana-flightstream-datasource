import React from 'react';
import { SelectableValue } from '@grafana/data';
import { InlineField, InlineFieldRow, Select } from '@grafana/ui';

// Mock data for live test channels
const liveTestDataChannels = [
  // {
  //   label: 'random-2s-stream',
  //   value: 'random-2s-stream',
  //   description: 'Random stream with points every 2s',
  // },
  // {
  //   label: 'random-flakey-stream',
  //   value: 'random-flakey-stream',
  //   description: 'Stream that returns data in random intervals',
  // },
  // {
  //   label: 'random-labeled-stream',
  //   value: 'random-labeled-stream',
  //   description: 'Value with moving labels',
  // },
  // {
  //   label: 'random-20Hz-stream',
  //   value: 'random-20Hz-stream',
  //   description: 'Random stream with points in 20Hz',
  // },
  {
    label: 'orcastream',
    value: 'orcastream',
    description: 'Orcastream flight',
  },
];

export function QueryEditor() {
  // Noop function: selection changes are ignored.
  const noop = (_: SelectableValue<string>) => {};

  // Use a static default selection.
  const selected = liveTestDataChannels[0];

  return (
    <InlineFieldRow>
      <InlineField label="Channel" labelWidth={14}>
        <Select
          width={32}
          onChange={noop}
          placeholder="Select channel"
          options={liveTestDataChannels}
          value={selected}
        />
      </InlineField>
    </InlineFieldRow>
  );
}