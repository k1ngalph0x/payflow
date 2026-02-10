import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:payflow/api/api_repo.dart';
import 'package:payflow/common/models/transaction_model.dart';
import 'package:payflow/common/utils/app_colors.dart';

class WalletController extends GetxController {
  final _apiRepo = ApiRepo();

  final RxBool isLoading = false.obs;
  final RxDouble balance = 0.0.obs;
  final RxList<TransactionModel> transactions = <TransactionModel>[].obs;

  @override
  void onInit() {
    super.onInit();
    //fetchBalance();
  }

  Future<void> fetchBalance() async {
    try {
      isLoading.value = true;

      final response = await _apiRepo.getWalletBalance();

      if (response['balance'] != null) {
        balance.value = (response['balance'] ?? 0.0).toDouble();
      }
    } catch (e) {
      _handleError(e);
    } finally {
      isLoading.value = false;
    }
  }

  Future<void> fetchTransactions({int limit = 20, int offset = 0}) async {
    try {
      isLoading.value = true;

      final response = await _apiRepo.getWalletTransactions(
        limit: limit,
        offset: offset,
      );

      if (response['transactions'] != null) {
        final List<dynamic> txnList = response['transactions'];
        transactions.value = txnList
            .map((json) => TransactionModel.fromJson(json))
            .toList();
      }
    } catch (e) {
      _handleError(e);
    } finally {
      isLoading.value = false;
    }
  }

  void _handleError(dynamic error) {
    String errorMessage = 'Failed to fetch wallet data';

    if (error.toString().contains('error')) {
      errorMessage = error.toString().split('error:').last.trim();
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
}
