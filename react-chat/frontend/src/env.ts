interface Env {
    SERVER_ADDRESS: string;
}

let serverAddress: string = '';

if (process.env.SERVER_ADDRESS) {
  serverAddress = process.env.SERVER_ADDRESS;
} else {
  throw new ReferenceError('process.env.SERVER_ADDRESS');
}

const env: Env = {
  SERVER_ADDRESS: serverAddress,
};
export default env;
