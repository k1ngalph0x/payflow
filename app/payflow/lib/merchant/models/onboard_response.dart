class OnboardResponse {
  final String merchantId;
  final String status;
  final String? message;

  OnboardResponse({
    required this.merchantId,
    required this.status,
    this.message,
  });

  factory OnboardResponse.fromJson(Map<String, dynamic> json) {
    return OnboardResponse(
      merchantId: json['merchant_id'] ?? '',
      status: json['status'] ?? 'PENDING',
      message: json['message'],
    );
  }
}
