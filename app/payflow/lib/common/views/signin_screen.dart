import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get.dart';
import 'package:payflow/common/controllers/auth_controller.dart';
import 'package:payflow/common/utils/app_colors.dart';
import 'package:payflow/common/utils/app_routes.dart';
import 'package:payflow/common/utils/validators.dart';
import 'package:payflow/common/widgets/custom_button.dart';
import 'package:payflow/common/widgets/custom_textfield.dart';

class SignInScreen extends StatefulWidget {
  const SignInScreen({super.key});

  @override
  State<SignInScreen> createState() => _SignInScreenState();
}

class _SignInScreenState extends State<SignInScreen> {
  final _authController = Get.find<AuthController>();
  final _formKey = GlobalKey<FormState>();
  final _emailController = TextEditingController();
  final _passwordController = TextEditingController();

  bool _isPasswordVisible = false;
  String _selectedRole = 'user';

  @override
  void initState() {
    super.initState();
    final args = Get.arguments;
    if (args != null && args['role'] != null) {
      _selectedRole = args['role'];
    }
  }

  @override
  void dispose() {
    _emailController.dispose();
    _passwordController.dispose();
    super.dispose();
  }

  void _handleSignIn() {
    if (_formKey.currentState!.validate()) {
      _authController.signIn(
        email: _emailController.text.trim(),
        password: _passwordController.text.trim(),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    final isUser = _selectedRole == 'user';
    final roleColor = isUser ? AppColors.primary : AppColors.success;

    return Scaffold(
      backgroundColor: Colors.white,
      body: SafeArea(
        child: SingleChildScrollView(
          padding: EdgeInsets.all(60.w),
          child: Form(
            key: _formKey,
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                IconButton(
                  icon: Icon(Icons.arrow_back, size: 60.sp),
                  onPressed: () => Get.back(),
                  padding: EdgeInsets.zero,
                  constraints: const BoxConstraints(),
                ),

                SizedBox(height: 48.h),

                Text(
                  'Welcome Back',
                  style: TextStyle(
                    fontSize: 72.sp,
                    fontWeight: FontWeight.bold,
                    color: AppColors.textPrimary,
                  ),
                ),

                SizedBox(height: 16.h),

                Container(
                  padding: EdgeInsets.symmetric(
                    horizontal: 24.w,
                    vertical: 12.h,
                  ),
                  decoration: BoxDecoration(
                    color: roleColor.withOpacity(0.1),
                    borderRadius: BorderRadius.circular(12.r),
                  ),
                  child: Row(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      Icon(
                        isUser ? Icons.person : Icons.store,
                        size: 28.sp,
                        color: roleColor,
                      ),
                      SizedBox(width: 8.w),
                      Text(
                        isUser ? 'User' : 'Merchant',
                        style: TextStyle(
                          fontSize: 28.sp,
                          fontWeight: FontWeight.w600,
                          color: roleColor,
                        ),
                      ),
                    ],
                  ),
                ),

                SizedBox(height: 96.h),

                CustomTextField(
                  controller: _emailController,
                  hintText: 'Email',
                  keyboardType: TextInputType.emailAddress,
                  prefixIcon: Icon(Icons.email_outlined, size: 48.sp),
                  validator: Validators.validateEmail,
                ),

                SizedBox(height: 32.h),

                CustomTextField(
                  controller: _passwordController,
                  hintText: 'Password',
                  obscureText: !_isPasswordVisible,
                  prefixIcon: Icon(Icons.lock_outline, size: 48.sp),
                  suffixIcon: IconButton(
                    icon: Icon(
                      _isPasswordVisible
                          ? Icons.visibility_outlined
                          : Icons.visibility_off_outlined,
                      size: 48.sp,
                    ),
                    onPressed: () => setState(
                      () => _isPasswordVisible = !_isPasswordVisible,
                    ),
                  ),
                  validator: Validators.validatePassword,
                ),

                SizedBox(height: 96.h),

                Obx(
                  () => CustomButton(
                    text: 'Sign In',
                    onPressed: _handleSignIn,
                    isLoading: _authController.isLoading.value,
                  ),
                ),

                SizedBox(height: 32.h),

                Center(
                  child: TextButton(
                    onPressed: () => Get.toNamed(
                      AppRoutes.signUp,
                      arguments: {'role': _selectedRole},
                    ),
                    child: RichText(
                      text: TextSpan(
                        style: TextStyle(fontSize: 28.sp),
                        children: [
                          TextSpan(
                            text: 'Don\'t have an account? ',
                            style: TextStyle(color: AppColors.textSecondary),
                          ),
                          TextSpan(
                            text: 'Sign Up',
                            style: TextStyle(
                              color: roleColor,
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                        ],
                      ),
                    ),
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
