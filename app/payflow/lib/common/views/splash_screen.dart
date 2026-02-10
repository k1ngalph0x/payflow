import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get.dart';
import 'package:payflow/common/controllers/auth_controller.dart';
import 'package:payflow/common/utils/app_colors.dart';
import 'package:payflow/common/utils/app_routes.dart';

class SplashScreen extends StatefulWidget {
  const SplashScreen({super.key});

  @override
  State<SplashScreen> createState() => _SplashScreenState();
}

class _SplashScreenState extends State<SplashScreen> {
  @override
  void initState() {
    super.initState();
    _navigate();
  }

  Future<void> _navigate() async {
    final authController = Get.find<AuthController>();

    await Future.delayed(const Duration(seconds: 2));

    if (authController.isAuthenticated.value && authController.isMerchant) {
      await authController.checkMerchantStatus();
    }

    if (authController.isAuthenticated.value) {
      if (authController.isUser) {
        Get.offAllNamed(AppRoutes.userHome);
      } else if (authController.isMerchant) {
        if (authController.isMerchantOnboarded.value) {
          Get.offAllNamed(AppRoutes.merchantHome);
        } else {
          Get.offAllNamed(AppRoutes.merchantOnboarding);
        }
      }
    } else {
      Get.offAllNamed(AppRoutes.roleSelection);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.primary,
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(
              Icons.account_balance_wallet_rounded,
              size: 200.sp,
              color: Colors.white,
            ),
            SizedBox(height: 32.h),
            Text(
              'PayFlow',
              style: TextStyle(
                fontSize: 84.sp,
                fontWeight: FontWeight.bold,
                color: Colors.white,
              ),
            ),
          ],
        ),
      ),
    );
  }
}
