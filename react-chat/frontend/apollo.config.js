/* eslint-disable @typescript-eslint/no-var-requires */
const path = require('path');

module.exports = {
  client: {
    includes: ['src/**/*'],
    service: {
      name: 'query',
      localSchemaFile: path.resolve(__dirname, '../schema.graphql'),
    },
  },
};
