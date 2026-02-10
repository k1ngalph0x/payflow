import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:payflow/api/api_repo.dart';
import 'package:payflow/common/models/auth_response.dart';
import 'package:payflow/common/utils/app_colors.dart';
import 'package:payflow/common/utils/app_routes.dart';
import 'package:payflow/common/utils/constants.dart';
import 'package:payflow/common/utils/shared_pref.dart';
import 'package:payflow/merchant/models/merchant_model.dart';
import 'package:payflow/merchant/models/onboard_response.dart';

class AuthController extends GetxController {
  final _apiRepo = ApiRepo();

  final RxBool isLoading = false.obs;
  final Rxn<String> authToken = Rxn<String>();
  final Rxn<String> userRole = Rxn<String>();
  final Rxn<String> userEmail = Rxn<String>();
  final Rxn<String> userId = Rxn<String>();
  final RxBool isAuthenticated = false.obs;

  final Rxn<MerchantModel> merchantData = Rxn<MerchantModel>();
  final RxBool isMerchantOnboarded = false.obs;
  final Rxn<String> merchantId = Rxn<String>();
  final Rxn<String> businessName = Rxn<String>();

  Map<String, String> get header => {
    "Content-Type": "application/json",
    "Accept": "application/json",
    "Authorization": "Bearer ${authToken.value}",
  };

  @override
  void onInit() {
    super.onInit();
    _initializeAuth();
  }

  Future<void> _initializeAuth() async {
    authToken.value = SharedPreferenceUtil.getString(Constants.keyAuthToken);
    userRole.value = SharedPreferenceUtil.getString(Constants.keyUserRole);
    userEmail.value = SharedPreferenceUtil.getString(Constants.keyUserEmail);
    userId.value = SharedPreferenceUtil.getString(Constants.keyUserId);

    String token = authToken.value = SharedPreferenceUtil.getString(
      Constants.keyAuthToken,
    );
    print('TOKEN $token');

    if (userRole.value == 'merchant') {
      merchantId.value = SharedPreferenceUtil.getString(
        Constants.keyMerchantId,
      );
      businessName.value = SharedPreferenceUtil.getString(
        Constants.keyBusinessName,
      );
      isMerchantOnboarded.value = SharedPreferenceUtil.getBool(
        Constants.keyIsOnboarded,
      );
    }

    if (authToken.value != null && authToken.value!.isNotEmpty) {
      isAuthenticated.value = true;
    }
  }

  Future<void> signIn({required String email, required String password}) async {
    try {
      isLoading.value = true;

      final response = await _apiRepo.signIn(email: email, password: password);
      final authResponse = AuthResponse.fromJson(response);

      await _saveUserData(
        token: authResponse.token,
        role: authResponse.role,
        email: email,
      );

      _showSuccessSnackbar(authResponse.message);

      await Future.delayed(const Duration(milliseconds: 500));
      await _navigateToHome(authResponse.role);
    } catch (e) {
      _handleError(e);
    } finally {
      isLoading.value = false;
    }
  }

  Future<void> signUp({
    required String email,
    required String password,
    required String role,
  }) async {
    try {
      isLoading.value = true;

      final response = await _apiRepo.signUp(
        email: email,
        password: password,
        role: role,
      );

      final authResponse = AuthResponse.fromJson(response);

      await _saveUserData(
        token: authResponse.token,
        role: authResponse.role,
        email: email,
      );

      _showSuccessSnackbar(authResponse.message);

      await Future.delayed(const Duration(milliseconds: 500));
      await _navigateToHome(authResponse.role);
    } catch (e) {
      _handleError(e);
    } finally {
      isLoading.value = false;
    }
  }

  Future<void> logout() async {
    try {
      await SharedPreferenceUtil.remove(Constants.keyAuthToken);
      await SharedPreferenceUtil.remove(Constants.keyUserRole);
      await SharedPreferenceUtil.remove(Constants.keyUserEmail);
      await SharedPreferenceUtil.remove(Constants.keyUserId);
      await SharedPreferenceUtil.remove(Constants.keyMerchantId);
      await SharedPreferenceUtil.remove(Constants.keyBusinessName);
      await SharedPreferenceUtil.remove(Constants.keyIsOnboarded);

      authToken.value = null;
      userRole.value = null;
      userEmail.value = null;
      userId.value = null;
      merchantId.value = null;
      businessName.value = null;
      merchantData.value = null;
      isAuthenticated.value = false;
      isMerchantOnboarded.value = false;

      Get.offAllNamed(AppRoutes.roleSelection);
    } catch (e) {
      print('Logout error: $e');
    }
  }

  Future<void> merchantOnboard({required String businessName}) async {
    try {
      isLoading.value = true;

      final response = await _apiRepo.merchantOnboard(
        businessName: businessName,
      );

      final onboardResponse = OnboardResponse.fromJson(response);

      await _saveMerchantData(
        merchantId: onboardResponse.merchantId,
        businessName: businessName,
        status: onboardResponse.status,
      );

      _showSuccessSnackbar('Merchant onboarded successfully!');

      await Future.delayed(const Duration(milliseconds: 500));
      Get.offAllNamed(AppRoutes.merchantHome);
    } catch (e) {
      _handleError(e);
    } finally {
      isLoading.value = false;
    }
  }

  Future<void> checkMerchantStatus() async {
    try {
      final response = await _apiRepo.getMerchantStatus();

      if (response['onboarded'] == true) {
        await _saveMerchantData(
          merchantId: response['merchant_id'],
          businessName: response['business_name'],
          status: response['status'],
        );
      }
    } catch (e) {
      print('Check merchant status error: $e');
    }
  }

  Future<void> _saveUserData({
    required String token,
    required String role,
    required String email,
    String? id,
  }) async {
    await SharedPreferenceUtil.setString(Constants.keyAuthToken, token);
    await SharedPreferenceUtil.setString(Constants.keyUserRole, role);
    await SharedPreferenceUtil.setString(Constants.keyUserEmail, email);

    if (id != null) {
      await SharedPreferenceUtil.setString(Constants.keyUserId, id);
    }

    authToken.value = token;
    userRole.value = role;
    userEmail.value = email;
    userId.value = id;
    isAuthenticated.value = true;
  }

  Future<void> _saveMerchantData({
    required String merchantId,
    required String businessName,
    required String status,
  }) async {
    await SharedPreferenceUtil.setString(Constants.keyMerchantId, merchantId);
    await SharedPreferenceUtil.setString(
      Constants.keyBusinessName,
      businessName,
    );
    await SharedPreferenceUtil.setBool(Constants.keyIsOnboarded, true);

    this.merchantId.value = merchantId;
    this.businessName.value = businessName;
    isMerchantOnboarded.value = true;

    merchantData.value = MerchantModel(
      merchantId: merchantId,
      userId: userId.value ?? '',
      businessName: businessName,
      status: status,
      createdAt: DateTime.now(),
    );
  }

  Future<void> _navigateToHome(String role) async {
    if (role == 'user') {
      Get.offAllNamed(AppRoutes.userHome);
    } else if (role == 'merchant') {
      await checkMerchantStatus();

      if (isMerchantOnboarded.value) {
        Get.offAllNamed(AppRoutes.merchantHome);
      } else {
        Get.offAllNamed(AppRoutes.merchantOnboarding);
      }
    }
  }

  void _handleError(dynamic error) {
    String errorMessage = 'An error occurred';

    if (error.toString().contains('error')) {
      errorMessage = error.toString().split('error:').last.trim();
    } else if (error.toString().contains('Email already exists')) {
      errorMessage = 'Email already exists';
    } else if (error.toString().contains('Invalid credentials')) {
      errorMessage = 'Invalid email or password';
    } else if (error.toString().contains('Merchant already exists')) {
      errorMessage = 'Merchant profile already exists';
    } else if (error.toString().contains('Only merchants can onboard')) {
      errorMessage = 'Only merchant accounts can onboard';
    }

    Get.snackbar(
      'Error',
      errorMessage,
      snackPosition: SnackPosition.BOTTOM,
      backgroundColor: AppColors.error,
      colorText: Colors.white,
      duration: const Duration(seconds: 3),
    );
  }

  void _showSuccessSnackbar(String message) {
    Get.snackbar(
      'Success',
      message,
      snackPosition: SnackPosition.BOTTOM,
      backgroundColor: AppColors.success,
      colorText: Colors.white,
      duration: const Duration(seconds: 2),
    );
  }

  bool get isUser => userRole.value == 'user';
  bool get isMerchant => userRole.value == 'merchant';
}
