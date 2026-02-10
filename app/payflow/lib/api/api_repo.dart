import 'package:payflow/api/api_config.dart';
import 'package:payflow/api/endpoints.dart';
import 'package:payflow/common/utils/constants.dart';
import 'package:uuid/uuid.dart';

class ApiRepo {
  final _apiConfig = ApiConfig();
  Future<Map<String, dynamic>> signIn({
    required String email,
    required String password,
  }) async {
    try {
      return await _apiConfig.postCall(
        baseUrl: Constants.authServiceUrl,
        endpoint: ApiEndpoints.signIn,
        body: {'email': email, 'password': password},
        isAuthHeader: false,
      );
    } catch (e) {
      rethrow;
    }
  }

  Future<Map<String, dynamic>> signUp({
    required String email,
    required String password,
    required String role,
  }) async {
    try {
      return await _apiConfig.postCall(
        baseUrl: Constants.authServiceUrl,
        endpoint: ApiEndpoints.signUp,
        body: {'email': email, 'password': password, 'role': role},
        isAuthHeader: false,
      );
    } catch (e) {
      rethrow;
    }
  }

  Future<Map<String, dynamic>> getProfile() async {
    try {
      return await _apiConfig.getCall(
        baseUrl: Constants.authServiceUrl,
        endpoint: ApiEndpoints.profile,
        isAuthHeader: true,
      );
    } catch (e) {
      rethrow;
    }
  }

  Future<Map<String, dynamic>> merchantOnboard({
    required String businessName,
  }) async {
    try {
      print('API: Merchant onboarding for $businessName');
      return await _apiConfig.postCall(
        baseUrl: Constants.merchantServiceUrl,
        endpoint: ApiEndpoints.merchantOnboard,
        body: {'business_name': businessName},
        isAuthHeader: true,
      );
    } catch (e) {
      print('API: Merchant onboard error - $e');
      rethrow;
    }
  }

  Future<Map<String, dynamic>> getMerchantProfile() async {
    try {
      print('API: Fetching merchant profile');
      return await _apiConfig.getCall(
        baseUrl: Constants.merchantServiceUrl,
        endpoint: ApiEndpoints.merchantProfile,
        isAuthHeader: true,
      );
    } catch (e) {
      print('API: Get merchant profile error - $e');
      rethrow;
    }
  }

  Future<Map<String, dynamic>> getMerchantStatus() async {
    try {
      print('API: Fetching merchant status');

      return await _apiConfig.getCall(
        baseUrl: Constants.merchantServiceUrl,
        endpoint: ApiEndpoints.merchantStatus,
        isAuthHeader: true,
      );
    } catch (e) {
      print('API: Get merchant status error - $e');
      rethrow;
    }
  }

  Future<Map<String, dynamic>> createPayment({
    required String merchantId,
    required double amount,
  }) async {
    try {
      final idempotencyKey = const Uuid().v4();
      print('API: Creating payment for merchant $merchantId, amount: $amount');
      return await _apiConfig.postCall(
        endpoint: ApiEndpoints.createPayment,
        baseUrl: Constants.paymentServiceUrl,
        body: {'merchant_id': merchantId, 'amount': amount},
        isAuthHeader: true,
        customHeaders: {'Idempotency-Key': idempotencyKey},
      );
    } catch (e) {
      print('API: Create payment error - $e');
      rethrow;
    }
  }

  Future<Map<String, dynamic>> getPaymentHistory({
    int limit = 20,
    int offset = 0,
  }) async {
    try {
      print('API: Fetching payment history');
      return await _apiConfig.getCall(
        baseUrl: Constants.paymentServiceUrl,
        endpoint: '${ApiEndpoints.paymentHistory}?limit=$limit&offset=$offset',
        isAuthHeader: true,
      );
    } catch (e) {
      print('API: Get payment history error - $e');
      rethrow;
    }
  }

  Future<Map<String, dynamic>> getPaymentStatus({
    required String reference,
  }) async {
    try {
      print('API: Fetching payment status for $reference');
      return await _apiConfig.getCall(
        baseUrl: Constants.paymentServiceUrl,
        endpoint: '${ApiEndpoints.paymentStatus}?reference=$reference',
        isAuthHeader: true,
      );
    } catch (e) {
      print('API: Get payment status error - $e');
      rethrow;
    }
  }

  Future<Map<String, dynamic>> getMerchantList() async {
    try {
      print('API: Fetching merchant list');
      return await _apiConfig.getCall(
        baseUrl: Constants.merchantServiceUrl,
        endpoint: ApiEndpoints.merchantList,
        isAuthHeader: true,
      );
    } catch (e) {
      print('API: Get merchant list error - $e');
      rethrow;
    }
  }

  Future<Map<String, dynamic>> getWalletBalance() async {
    try {
      print('API: Fetching wallet balance');
      return await _apiConfig.getCall(
        baseUrl: Constants.walletServiceUrl,
        endpoint: ApiEndpoints.walletBalance,
        isAuthHeader: true,
      );
    } catch (e) {
      print('API: Get wallet balance error - $e');
      rethrow;
    }
  }

  Future<Map<String, dynamic>> getWalletTransactions({
    int limit = 20,
    int offset = 0,
  }) async {
    try {
      print('API: Fetching wallet transactions');
      return await _apiConfig.getCall(
        baseUrl: Constants.walletServiceUrl,
        endpoint:
            '${ApiEndpoints.walletTransactions}?limit=$limit&offset=$offset',
        isAuthHeader: true,
      );
    } catch (e) {
      print('API: Get wallet transactions error - $e');
      rethrow;
    }
  }
}
