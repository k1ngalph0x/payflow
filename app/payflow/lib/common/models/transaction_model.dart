class TransactionModel {
  final String id;
  final String type;
  final double amount;
  final String reference;
  final String status;
  final DateTime createdAt;

  TransactionModel({
    required this.id,
    required this.type,
    required this.amount,
    required this.reference,
    required this.status,
    required this.createdAt,
  });

  factory TransactionModel.fromJson(Map<String, dynamic> json) {
    return TransactionModel(
      id: json['id'] ?? '',
      type: json['type'] ?? '',
      amount: (json['amount'] ?? 0.0).toDouble(),
      reference: json['reference'] ?? '',
      status: json['status'] ?? '',
      createdAt: json['created_at'] != null
          ? DateTime.parse(json['created_at'])
          : DateTime.now(),
    );
  }

  bool get isDebit => type == 'DEBIT';
  bool get isCredit => type == 'CREDIT';
}
