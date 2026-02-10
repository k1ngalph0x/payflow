import 'dart:convert';

import 'package:get/get.dart';
import 'package:http/http.dart' as http;
import 'package:payflow/common/controllers/auth_controller.dart';
import 'package:payflow/common/utils/constants.dart';

class ApiConfig {
  //static const String baseUrl = Constants.baseUrl;
  static const Duration timeoutDuration = Duration(seconds: 30);

  final Map<String, String> _defaultHeaders = {
    "Content-Type": "application/json",
    "Accept": "application/json",
  };

  dynamic _handleResponse(http.Response response) {
    if (response.statusCode >= 200 && response.statusCode < 300) {
      return jsonDecode(response.body);
    } else {
      throw jsonDecode(response.body);
    }
  }

  Map<String, String> _getHeaders({
    required bool isAuthHeader,
    Map<String, String>? customHeaders,
  }) {
    Map<String, String> headers = Map.from(_defaultHeaders);

    if (isAuthHeader) {
      headers.addAll(Get.find<AuthController>().header);
    }

    if (customHeaders != null) {
      headers.addAll(customHeaders);
    }

    return headers;
  }

  Future<dynamic> getCall({
    required String baseUrl,
    required String endpoint,
    required bool isAuthHeader,
    Map<String, String>? customHeaders,
  }) async {
    try {
      final response = await http
          .get(
            Uri.parse("$baseUrl$endpoint"),
            headers: _getHeaders(
              isAuthHeader: isAuthHeader,
              customHeaders: customHeaders,
            ),
          )
          .timeout(timeoutDuration);

      return _handleResponse(response);
    } catch (e) {
      rethrow;
    }
  }

  Future<dynamic> postCall({
    required String baseUrl,
    required String endpoint,
    required Map<String, dynamic> body,
    required bool isAuthHeader,
    Map<String, String>? customHeaders,
  }) async {
    try {
      final response = await http
          .post(
            Uri.parse("$baseUrl$endpoint"),
            headers: _getHeaders(
              isAuthHeader: isAuthHeader,
              customHeaders: customHeaders,
            ),
            body: jsonEncode(body),
          )
          .timeout(timeoutDuration);

      return _handleResponse(response);
    } catch (e) {
      rethrow;
    }
  }

  Future<dynamic> putCall({
    required String baseUrl,
    required String endpoint,
    required Map<String, dynamic> body,
    required bool isAuthHeader,
    Map<String, String>? customHeaders,
  }) async {
    try {
      final response = await http
          .put(
            Uri.parse("$baseUrl$endpoint"),
            headers: _getHeaders(
              isAuthHeader: isAuthHeader,
              customHeaders: customHeaders,
            ),
            body: jsonEncode(body),
          )
          .timeout(timeoutDuration);

      return _handleResponse(response);
    } catch (e) {
      rethrow;
    }
  }

  Future<dynamic> patchCall({
    required String baseUrl,
    required String endpoint,
    required Map<String, dynamic> body,
    required bool isAuthHeader,
    Map<String, String>? customHeaders,
  }) async {
    try {
      final response = await http
          .patch(
            Uri.parse("$baseUrl$endpoint"),
            headers: _getHeaders(
              isAuthHeader: isAuthHeader,
              customHeaders: customHeaders,
            ),
            body: jsonEncode(body),
          )
          .timeout(timeoutDuration);

      return _handleResponse(response);
    } catch (e) {
      rethrow;
    }
  }

  Future<dynamic> deleteCall({
    required String baseUrl,
    required String endpoint,
    required bool isAuthHeader,
    Map<String, dynamic>? body,
    Map<String, String>? customHeaders,
  }) async {
    try {
      final response = await http
          .delete(
            Uri.parse("$baseUrl$endpoint"),
            headers: _getHeaders(
              isAuthHeader: isAuthHeader,
              customHeaders: customHeaders,
            ),
            body: body != null ? jsonEncode(body) : null,
          )
          .timeout(timeoutDuration);

      return _handleResponse(response);
    } catch (e) {
      rethrow;
    }
  }
}
