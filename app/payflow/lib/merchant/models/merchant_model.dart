class MerchantModel {
  final String merchantId;
  final String userId;
  final String businessName;
  final String status;
  final DateTime createdAt;

  MerchantModel({
    required this.merchantId,
    required this.userId,
    required this.businessName,
    required this.status,
    required this.createdAt,
  });

  factory MerchantModel.fromJson(Map<String, dynamic> json) {
    return MerchantModel(
      merchantId: json['merchant_id'] ?? json['id'] ?? '',
      userId: json['user_id'] ?? '',
      businessName: json['business_name'] ?? '',
      status: json['status'] ?? 'PENDING',
      createdAt: json['created_at'] != null
          ? DateTime.parse(json['created_at'])
          : DateTime.now(),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'merchant_id': merchantId,
      'user_id': userId,
      'business_name': businessName,
      'status': status,
      'created_at': createdAt.toIso8601String(),
    };
  }

  bool get isActive => status == 'ACTIVE';
  bool get isPending => status == 'PENDING';
  bool get isRejected => status == 'REJECTED';
}
