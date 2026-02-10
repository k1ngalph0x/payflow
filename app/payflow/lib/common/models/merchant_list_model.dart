class MerchantListModel {
  final String id;
  final String userId;
  final String businessName;
  final String status;
  final DateTime createdAt;

  MerchantListModel({
    required this.id,
    required this.userId,
    required this.businessName,
    required this.status,
    required this.createdAt,
  });

  factory MerchantListModel.fromJson(Map<String, dynamic> json) {
    return MerchantListModel(
      id: json['id'] ?? '',
      userId: json['user_id'] ?? '',
      businessName: json['business_name'] ?? '',
      status: json['status'] ?? 'PENDING',
      createdAt: json['created_at'] != null
          ? DateTime.parse(json['created_at'])
          : DateTime.now(),
    );
  }
}
