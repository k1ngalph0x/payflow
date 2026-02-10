import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:payflow/api/api_repo.dart';
import 'package:payflow/common/controllers/auth_controller.dart';
import 'package:payflow/common/models/create_payment_response.dart';
import 'package:payflow/common/models/payment_model.dart';
import 'package:payflow/common/utils/app_colors.dart';

class PaymentController extends GetxController {
  final _apiRepo = ApiRepo();
  final authController = Get.find<AuthController>();

  final RxBool isLoading = false.obs;
  final RxBool isCreatingPayment = false.obs;
  final RxList<PaymentModel> payments = <PaymentModel>[].obs;
  final Rxn<PaymentModel> currentPayment = Rxn<PaymentModel>();

  Future<CreatePaymentResponse?> createPayment({
    required String merchantId,
    required double amount,
  }) async {
    try {
      isCreatingPayment.value = true;

      final response = await _apiRepo.createPayment(
        merchantId: merchantId,
        amount: amount,
      );

      final paymentResponse = CreatePaymentResponse.fromJson(response);

      _showSuccessSnackbar('Payment initiated successfully!');

      return paymentResponse;
    } catch (e) {
      _handleError(e);
      return null;
    } finally {
      isCreatingPayment.value = false;
    }
  }

  Future<void> getPaymentHistory({int limit = 20, int offset = 0}) async {
    try {
      isLoading.value = true;

      final response = await _apiRepo.getPaymentHistory(
        limit: limit,
        offset: offset,
      );

      if (response['payments'] != null) {
        final List<dynamic> paymentList = response['payments'];
        payments.value = paymentList
            .map((json) => PaymentModel.fromJson(json))
            .toList();
      }
    } catch (e) {
      _handleError(e);
    } finally {
      isLoading.value = false;
    }
  }

  Future<PaymentModel?> checkPaymentStatus({required String reference}) async {
    try {
      final response = await _apiRepo.getPaymentStatus(reference: reference);

      if (response['payment'] != null) {
        return PaymentModel.fromJson(response['payment']);
      }
      return null;
    } catch (e) {
      print('Check payment status error: $e');
      return null;
    }
  }

  Future<void> pollPaymentStatus({
    required String reference,
    required Function(PaymentModel) onStatusUpdate,
    int maxAttempts = 30,
    Duration interval = const Duration(seconds: 2),
  }) async {
    int attempts = 0;

    while (attempts < maxAttempts) {
      final payment = await checkPaymentStatus(reference: reference);

      if (payment != null) {
        onStatusUpdate(payment);

        if (payment.isFundsCaptured || payment.isFailed) {
          break;
        }
      }

      await Future.delayed(interval);
      attempts++;
    }
  }

  void _handleError(dynamic error) {
    String errorMessage = 'An error occurred';

    if (error.toString().contains('error')) {
      errorMessage = error.toString().split('error:').last.trim();
    } else if (error.toString().contains('Insufficient')) {
      errorMessage = 'Insufficient funds in wallet';
    } else if (error.toString().contains('Invalid')) {
      errorMessage = 'Invalid payment details';
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
}
