class CreatePaymentResponse {
  final String reference;
  final String status;
  final String? message;

  CreatePaymentResponse({
    required this.reference,
    required this.status,
    this.message,
  });

  factory CreatePaymentResponse.fromJson(Map<String, dynamic> json) {
    return CreatePaymentResponse(
      reference: json['reference'] ?? '',
      status: json['status'] ?? 'CREATED',
      message: json['message'],
    );
  }
}
