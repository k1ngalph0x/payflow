class AuthResponse {
  final String message;
  final String token;
  final String role;

  AuthResponse({
    required this.message,
    required this.token,
    required this.role,
  });

  factory AuthResponse.fromJson(Map<String, dynamic> json) {
    return AuthResponse(
      message: json['message'] ?? '',
      token: json['token'] ?? '',
      role: json['role'] ?? '',
    );
  }

  Map<String, dynamic> toJson() {
    return {'message': message, 'token': token, 'role': role};
  }
}
