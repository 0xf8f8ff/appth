/**
 * @fileoverview gRPC-Web generated client stub for main
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


/* eslint-disable */
// @ts-nocheck



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.main = require('./appth_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.main.AuthClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options.format = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.main.AuthPromiseClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options.format = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.main.PingRequest,
 *   !proto.main.PongResponse>}
 */
const methodDescriptor_Auth_Ping = new grpc.web.MethodDescriptor(
  '/main.Auth/Ping',
  grpc.web.MethodType.UNARY,
  proto.main.PingRequest,
  proto.main.PongResponse,
  /**
   * @param {!proto.main.PingRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.main.PongResponse.deserializeBinary
);


/**
 * @param {!proto.main.PingRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.main.PongResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.main.PongResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.main.AuthClient.prototype.ping =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/main.Auth/Ping',
      request,
      metadata || {},
      methodDescriptor_Auth_Ping,
      callback);
};


/**
 * @param {!proto.main.PingRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.main.PongResponse>}
 *     Promise that resolves to the response
 */
proto.main.AuthPromiseClient.prototype.ping =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/main.Auth/Ping',
      request,
      metadata || {},
      methodDescriptor_Auth_Ping);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.main.UserRequest,
 *   !proto.main.User>}
 */
const methodDescriptor_Auth_GetUser = new grpc.web.MethodDescriptor(
  '/main.Auth/GetUser',
  grpc.web.MethodType.UNARY,
  proto.main.UserRequest,
  proto.main.User,
  /**
   * @param {!proto.main.UserRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.main.User.deserializeBinary
);


/**
 * @param {!proto.main.UserRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.main.User)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.main.User>|undefined}
 *     The XHR Node Readable Stream
 */
proto.main.AuthClient.prototype.getUser =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/main.Auth/GetUser',
      request,
      metadata || {},
      methodDescriptor_Auth_GetUser,
      callback);
};


/**
 * @param {!proto.main.UserRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.main.User>}
 *     Promise that resolves to the response
 */
proto.main.AuthPromiseClient.prototype.getUser =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/main.Auth/GetUser',
      request,
      metadata || {},
      methodDescriptor_Auth_GetUser);
};


module.exports = proto.main;

