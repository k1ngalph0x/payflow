import 'dart:ui';

import 'package:payflow/common/utils/app_colors.dart';

class PaymentModel {
  final String id;
  final String userId;
  final String merchantId;
  final double amount;
  final String status;
  final String reference;
  final DateTime createdAt;

  PaymentModel({
    required this.id,
    required this.userId,
    required this.merchantId,
    required this.amount,
    required this.status,
    required this.reference,
    required this.createdAt,
  });

  factory PaymentModel.fromJson(Map<String, dynamic> json) {
    return PaymentModel(
      id: json['id'] ?? '',
      userId: json['user_id'] ?? '',
      merchantId: json['merchant_id'] ?? '',
      amount: (json['amount'] ?? 0.0).toDouble(),
      status: json['status'] ?? 'PENDING',
      reference: json['reference'] ?? '',
      createdAt: json['created_at'] != null
          ? DateTime.parse(json['created_at'])
          : DateTime.now(),
    );
  }

  bool get isCreated => status == 'CREATED';
  bool get isProcessing => status == 'PROCESSING';
  bool get isFundsCaptured => status == 'FUNDS_CAPTURED';
  bool get isFailed => status == 'FAILED';
  bool get isSuccess => status == 'FUNDS_CAPTURED';

  Color getStatusColor() {
    switch (status) {
      case 'CREATED':
        return AppColors.warning;
      case 'PROCESSING':
        return AppColors.info;
      case 'FUNDS_CAPTURED':
        return AppColors.success;
      case 'FAILED':
        return AppColors.error;
      default:
        return AppColors.textSecondary;
    }
  }

  String getStatusText() {
    switch (status) {
      case 'CREATED':
        return 'Created';
      case 'PROCESSING':
        return 'Processing';
      case 'FUNDS_CAPTURED':
        return 'Completed';
      case 'FAILED':
        return 'Failed';
      default:
        return status;
    }
  }
}
