const abi = require('ethereumjs-abi');

export default function encodeCall(name, args, values) {
  const methodId = abi.methodID(name, args).toString('hex');
  const params = abi.rawEncode(args, values).toString('hex');
  return '0x' + methodId + params;
};
