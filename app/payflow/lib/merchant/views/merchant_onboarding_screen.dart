import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get.dart';
import 'package:payflow/common/controllers/auth_controller.dart';
import 'package:payflow/common/utils/app_colors.dart';
import 'package:payflow/common/widgets/custom_button.dart';
import 'package:payflow/common/widgets/custom_textfield.dart';

class MerchantOnboardingScreen extends StatefulWidget {
  const MerchantOnboardingScreen({super.key});

  @override
  State<MerchantOnboardingScreen> createState() =>
      _MerchantOnboardingScreenState();
}

class _MerchantOnboardingScreenState extends State<MerchantOnboardingScreen> {
  final _authController = Get.find<AuthController>();
  final _formKey = GlobalKey<FormState>();
  final _businessNameController = TextEditingController();

  @override
  void dispose() {
    _businessNameController.dispose();
    super.dispose();
  }

  void _handleOnboard() {
    if (_formKey.currentState!.validate()) {
      _authController.merchantOnboard(
        businessName: _businessNameController.text.trim(),
      );
    }
  }

  String? _validateBusinessName(String? value) {
    if (value == null || value.isEmpty) {
      return 'Business name is required';
    }
    if (value.length < 3) {
      return 'Business name must be at least 3 characters';
    }
    return null;
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.white,
      body: SafeArea(
        child: SingleChildScrollView(
          padding: EdgeInsets.symmetric(horizontal: 60.w),
          child: Form(
            key: _formKey,
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                SizedBox(height: 96.h),

                Center(
                  child: Container(
                    padding: EdgeInsets.all(72.w),
                    decoration: BoxDecoration(
                      color: AppColors.success.withOpacity(0.1),
                      borderRadius: BorderRadius.circular(48.r),
                    ),
                    child: Icon(
                      Icons.store_rounded,
                      size: 240.sp,
                      color: AppColors.success,
                    ),
                  ),
                ),

                SizedBox(height: 72.h),

                Text(
                  'Complete Your\nMerchant Profile',
                  style: TextStyle(
                    fontSize: 84.sp,
                    fontWeight: FontWeight.bold,
                    color: AppColors.textPrimary,
                    height: 1.2,
                  ),
                ),

                SizedBox(height: 24.h),

                Text(
                  'Tell us about your business to start accepting payments',
                  style: TextStyle(
                    fontSize: 36.sp,
                    color: AppColors.textSecondary,
                    height: 1.4,
                  ),
                ),

                SizedBox(height: 96.h),

                Text(
                  'Business Name',
                  style: TextStyle(
                    fontSize: 32.sp,
                    fontWeight: FontWeight.w600,
                    color: AppColors.textPrimary,
                  ),
                ),
                SizedBox(height: 24.h),

                CustomTextField(
                  controller: _businessNameController,
                  hintText: 'Enter your business name',
                  prefixIcon: Icon(
                    Icons.business_rounded,
                    size: 48.sp,
                    color: AppColors.textSecondary,
                  ),
                  validator: _validateBusinessName,
                ),

                SizedBox(height: 48.h),

                Container(
                  padding: EdgeInsets.all(48.w),
                  decoration: BoxDecoration(
                    color: AppColors.info.withOpacity(0.1),
                    borderRadius: BorderRadius.circular(24.r),
                    border: Border.all(
                      color: AppColors.info.withOpacity(0.3),
                      width: 1,
                    ),
                  ),
                  child: Row(
                    children: [
                      Icon(
                        Icons.info_outline_rounded,
                        size: 48.sp,
                        color: AppColors.info,
                      ),
                      SizedBox(width: 32.w),
                      Expanded(
                        child: Text(
                          'Your business name will be visible to customers when they make payments',
                          style: TextStyle(
                            fontSize: 28.sp,
                            color: AppColors.textSecondary,
                            height: 1.4,
                          ),
                        ),
                      ),
                    ],
                  ),
                ),

                SizedBox(height: 96.h),

                Obx(
                  () => CustomButton(
                    text: 'Complete Onboarding',
                    onPressed: _handleOnboard,
                    isLoading: _authController.isLoading.value,
                    backgroundColor: AppColors.success,
                  ),
                ),

                SizedBox(height: 48.h),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
