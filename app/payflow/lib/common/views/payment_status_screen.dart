import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get.dart';
import 'package:payflow/common/controllers/payment_controller.dart';
import 'package:payflow/common/models/payment_model.dart';
import 'package:payflow/common/utils/app_colors.dart';
import 'package:payflow/common/utils/app_routes.dart';
import 'package:payflow/common/widgets/custom_button.dart';

class PaymentStatusScreen extends StatefulWidget {
  final String reference;

  const PaymentStatusScreen({super.key, required this.reference});

  @override
  State<PaymentStatusScreen> createState() => _PaymentStatusScreenState();
}

class _PaymentStatusScreenState extends State<PaymentStatusScreen> {
  final _paymentController = Get.put(PaymentController());
  final Rxn<PaymentModel> paymentStatus = Rxn<PaymentModel>();

  @override
  void initState() {
    super.initState();
    _startPolling();
  }

  void _startPolling() {
    _paymentController.pollPaymentStatus(
      reference: widget.reference,
      onStatusUpdate: (payment) {
        setState(() {
          paymentStatus.value = payment;
        });
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.white,
      appBar: AppBar(
        backgroundColor: Colors.white,
        elevation: 0,
        leading: Container(),
        title: Text(
          'Payment Status',
          style: TextStyle(
            fontSize: 42.sp,
            fontWeight: FontWeight.bold,
            color: AppColors.textPrimary,
          ),
        ),
      ),
      body: Obx(() {
        final payment = paymentStatus.value;

        return Center(
          child: Padding(
            padding: EdgeInsets.all(60.w),
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Container(
                  padding: EdgeInsets.all(72.w),
                  decoration: BoxDecoration(
                    color: payment != null
                        ? payment.getStatusColor().withOpacity(0.1)
                        : AppColors.warning.withOpacity(0.1),
                    borderRadius: BorderRadius.circular(48.r),
                  ),
                  child: Icon(
                    payment != null && payment.isFundsCaptured
                        ? Icons.check_circle_rounded
                        : payment != null && payment.isFailed
                        ? Icons.error_rounded
                        : Icons.hourglass_empty_rounded,
                    size: 240.sp,
                    color: payment?.getStatusColor() ?? AppColors.warning,
                  ),
                ),

                SizedBox(height: 72.h),

                Text(
                  payment?.getStatusText() ?? 'Processing...',
                  style: TextStyle(
                    fontSize: 72.sp,
                    fontWeight: FontWeight.bold,
                    color: AppColors.textPrimary,
                  ),
                ),

                SizedBox(height: 24.h),

                Container(
                  padding: EdgeInsets.symmetric(
                    horizontal: 48.w,
                    vertical: 24.h,
                  ),
                  decoration: BoxDecoration(
                    color: AppColors.backgroundLight,
                    borderRadius: BorderRadius.circular(16.r),
                  ),
                  child: Text(
                    'Ref: ${widget.reference.substring(0, 8)}...',
                    style: TextStyle(
                      fontSize: 28.sp,
                      color: AppColors.textSecondary,
                      fontFamily: 'monospace',
                    ),
                  ),
                ),

                SizedBox(height: 48.h),

                if (payment != null)
                  Text(
                    '₹${payment.amount.toStringAsFixed(2)}',
                    style: TextStyle(
                      fontSize: 96.sp,
                      fontWeight: FontWeight.bold,
                      color: AppColors.textPrimary,
                    ),
                  ),

                const Spacer(),

                if (payment != null &&
                    (payment.isFundsCaptured || payment.isFailed))
                  CustomButton(
                    text: 'Done',
                    onPressed: () => Get.offAllNamed(AppRoutes.userHome),
                  ),

                SizedBox(height: 48.h),
              ],
            ),
          ),
        );
      }),
    );
  }
}
