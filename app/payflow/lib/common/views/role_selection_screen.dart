import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get.dart';
import 'package:payflow/common/utils/app_colors.dart';
import 'package:payflow/common/utils/app_routes.dart';

class RoleSelectionScreen extends StatelessWidget {
  const RoleSelectionScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.white,
      body: SafeArea(
        child: Padding(
          padding: EdgeInsets.all(60.w),
          child: Column(
            children: [
              const Spacer(),

              Icon(
                Icons.account_balance_wallet_rounded,
                size: 180.sp,
                color: AppColors.primary,
              ),
              SizedBox(height: 32.h),
              Text(
                'PayFlow',
                style: TextStyle(
                  fontSize: 84.sp,
                  fontWeight: FontWeight.bold,
                  color: AppColors.textPrimary,
                ),
              ),
              SizedBox(height: 16.h),
              Text(
                'Choose your account type',
                style: TextStyle(
                  fontSize: 32.sp,
                  color: AppColors.textSecondary,
                ),
              ),

              const Spacer(),

              _RoleCard(
                icon: Icons.person_rounded,
                title: 'User',
                subtitle: 'Send & receive payments',
                color: AppColors.primary,
                onTap: () =>
                    Get.toNamed(AppRoutes.signIn, arguments: {'role': 'user'}),
              ),

              SizedBox(height: 32.h),

              _RoleCard(
                icon: Icons.store_rounded,
                title: 'Merchant',
                subtitle: 'Accept payments from customers',
                color: AppColors.success,
                onTap: () => Get.toNamed(
                  AppRoutes.signIn,
                  arguments: {'role': 'merchant'},
                ),
              ),

              const Spacer(),
            ],
          ),
        ),
      ),
    );
  }
}

class _RoleCard extends StatelessWidget {
  final IconData icon;
  final String title;
  final String subtitle;
  final Color color;
  final VoidCallback onTap;

  const _RoleCard({
    required this.icon,
    required this.title,
    required this.subtitle,
    required this.color,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return Material(
      color: Colors.transparent,
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(24.r),
        child: Container(
          padding: EdgeInsets.all(48.w),
          decoration: BoxDecoration(
            color: color.withOpacity(0.05),
            borderRadius: BorderRadius.circular(24.r),
            border: Border.all(color: color.withOpacity(0.2), width: 2),
          ),
          child: Row(
            children: [
              Container(
                padding: EdgeInsets.all(32.w),
                decoration: BoxDecoration(
                  color: color.withOpacity(0.1),
                  borderRadius: BorderRadius.circular(16.r),
                ),
                child: Icon(icon, size: 72.sp, color: color),
              ),
              SizedBox(width: 32.w),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      title,
                      style: TextStyle(
                        fontSize: 42.sp,
                        fontWeight: FontWeight.bold,
                        color: AppColors.textPrimary,
                      ),
                    ),
                    SizedBox(height: 4.h),
                    Text(
                      subtitle,
                      style: TextStyle(
                        fontSize: 28.sp,
                        color: AppColors.textSecondary,
                      ),
                    ),
                  ],
                ),
              ),
              Icon(Icons.arrow_forward_ios, size: 36.sp, color: color),
            ],
          ),
        ),
      ),
    );
  }
}
