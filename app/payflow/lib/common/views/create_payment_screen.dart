import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get.dart';
import 'package:payflow/common/controllers/payment_controller.dart';
import 'package:payflow/common/utils/app_colors.dart';
import 'package:payflow/common/views/payment_status_screen.dart';
import 'package:payflow/common/widgets/custom_button.dart';
import 'package:payflow/common/widgets/custom_textfield.dart';
import 'package:payflow/merchant/controller/merchant_controller.dart';

class CreatePaymentScreen extends StatefulWidget {
  const CreatePaymentScreen({super.key});

  @override
  State<CreatePaymentScreen> createState() => _CreatePaymentScreenState();
}

class _CreatePaymentScreenState extends State<CreatePaymentScreen> {
  final _paymentController = Get.put(PaymentController());
  final _merchantController = Get.put(MerchantController());
  final _formKey = GlobalKey<FormState>();
  final _amountController = TextEditingController();

  @override
  void dispose() {
    _amountController.dispose();
    super.dispose();
  }

  String? _validateAmount(String? value) {
    if (value == null || value.isEmpty) {
      return 'Amount is required';
    }
    final amount = double.tryParse(value);
    if (amount == null || amount <= 0) {
      return 'Enter a valid amount';
    }
    return null;
  }

  String? _validateMerchant() {
    if (_merchantController.selectedMerchant.value == null) {
      return 'Please select a merchant';
    }
    return null;
  }

  Future<void> _handleCreatePayment() async {
    final merchantError = _validateMerchant();
    if (merchantError != null) {
      Get.snackbar(
        'Error',
        merchantError,
        snackPosition: SnackPosition.BOTTOM,
        backgroundColor: AppColors.error,
        colorText: Colors.white,
      );
      return;
    }

    if (_formKey.currentState!.validate()) {
      final response = await _paymentController.createPayment(
        merchantId: _merchantController.selectedMerchant.value!.id,
        amount: double.parse(_amountController.text.trim()),
      );

      if (response != null) {
        Get.to(() => PaymentStatusScreen(reference: response.reference));
      }
    }
  }

  void _showMerchantSelector() {
    showModalBottomSheet(
      context: context,
      backgroundColor: Colors.transparent,
      isScrollControlled: true,
      builder: (context) => _buildMerchantBottomSheet(),
    );
  }

  Widget _buildMerchantBottomSheet() {
    return Container(
      height: MediaQuery.of(context).size.height * 0.7,
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.vertical(top: Radius.circular(32.r)),
      ),
      child: Column(
        children: [
          Container(
            margin: EdgeInsets.only(top: 24.h),
            width: 100.w,
            height: 8.h,
            decoration: BoxDecoration(
              color: Colors.grey[300],
              borderRadius: BorderRadius.circular(4.r),
            ),
          ),

          Padding(
            padding: EdgeInsets.all(60.w),
            child: Row(
              children: [
                Text(
                  'Select Merchant',
                  style: TextStyle(
                    fontSize: 48.sp,
                    fontWeight: FontWeight.bold,
                    color: AppColors.textPrimary,
                  ),
                ),
                const Spacer(),
                IconButton(
                  icon: Icon(
                    Icons.close_rounded,
                    size: 56.sp,
                    color: AppColors.textPrimary,
                  ),
                  onPressed: () => Get.back(),
                ),
              ],
            ),
          ),

          Expanded(
            child: Obx(() {
              if (_merchantController.isLoading.value) {
                return Center(
                  child: CircularProgressIndicator(color: AppColors.primary),
                );
              }

              if (_merchantController.merchants.isEmpty) {
                return Center(
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Icon(
                        Icons.store_outlined,
                        size: 180.sp,
                        color: Colors.grey[300],
                      ),
                      SizedBox(height: 32.h),
                      Text(
                        'No merchants available',
                        style: TextStyle(
                          fontSize: 36.sp,
                          color: AppColors.textSecondary,
                        ),
                      ),
                    ],
                  ),
                );
              }

              return ListView.builder(
                padding: EdgeInsets.symmetric(horizontal: 60.w),
                itemCount: _merchantController.merchants.length,
                itemBuilder: (context, index) {
                  final merchant = _merchantController.merchants[index];
                  final isSelected =
                      _merchantController.selectedMerchant.value?.id ==
                      merchant.id;

                  return GestureDetector(
                    onTap: () {
                      _merchantController.selectMerchant(merchant);
                      Get.back();
                    },
                    child: Container(
                      margin: EdgeInsets.only(bottom: 32.h),
                      padding: EdgeInsets.all(48.w),
                      decoration: BoxDecoration(
                        color: isSelected
                            ? AppColors.primary.withOpacity(0.1)
                            : AppColors.backgroundLight,
                        borderRadius: BorderRadius.circular(24.r),
                        border: Border.all(
                          color: isSelected
                              ? AppColors.primary
                              : Colors.transparent,
                          width: 2,
                        ),
                      ),
                      child: Row(
                        children: [
                          Container(
                            padding: EdgeInsets.all(36.w),
                            decoration: BoxDecoration(
                              color: isSelected
                                  ? AppColors.primary.withOpacity(0.2)
                                  : Colors.grey[200],
                              borderRadius: BorderRadius.circular(16.r),
                            ),
                            child: Icon(
                              Icons.store_rounded,
                              size: 56.sp,
                              color: isSelected
                                  ? AppColors.primary
                                  : AppColors.textSecondary,
                            ),
                          ),
                          SizedBox(width: 32.w),
                          Expanded(
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                Text(
                                  merchant.businessName,
                                  style: TextStyle(
                                    fontSize: 36.sp,
                                    fontWeight: FontWeight.w600,
                                    color: AppColors.textPrimary,
                                  ),
                                ),
                                SizedBox(height: 8.h),
                                Text(
                                  'ID: ${merchant.id.substring(0, 8)}...',
                                  style: TextStyle(
                                    fontSize: 26.sp,
                                    color: AppColors.textSecondary,
                                    fontFamily: 'monospace',
                                  ),
                                ),
                              ],
                            ),
                          ),
                          if (isSelected)
                            Icon(
                              Icons.check_circle_rounded,
                              size: 56.sp,
                              color: AppColors.primary,
                            ),
                        ],
                      ),
                    ),
                  );
                },
              );
            }),
          ),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.white,
      appBar: AppBar(
        backgroundColor: Colors.white,
        elevation: 0,
        leading: IconButton(
          icon: Icon(
            Icons.arrow_back_ios_rounded,
            color: AppColors.textPrimary,
            size: 48.sp,
          ),
          onPressed: () => Get.back(),
        ),
        title: Text(
          'Make Payment',
          style: TextStyle(
            fontSize: 42.sp,
            fontWeight: FontWeight.bold,
            color: AppColors.textPrimary,
          ),
        ),
      ),
      body: SingleChildScrollView(
        padding: EdgeInsets.all(60.w),
        child: Form(
          key: _formKey,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              SizedBox(height: 48.h),

              Center(
                child: Container(
                  padding: EdgeInsets.all(72.w),
                  decoration: BoxDecoration(
                    color: AppColors.primary.withOpacity(0.1),
                    borderRadius: BorderRadius.circular(48.r),
                  ),
                  child: Icon(
                    Icons.payment_rounded,
                    size: 180.sp,
                    color: AppColors.primary,
                  ),
                ),
              ),

              SizedBox(height: 72.h),

              Text(
                'Pay To',
                style: TextStyle(
                  fontSize: 32.sp,
                  fontWeight: FontWeight.w600,
                  color: AppColors.textPrimary,
                ),
              ),
              SizedBox(height: 24.h),
              Obx(() {
                final selectedMerchant =
                    _merchantController.selectedMerchant.value;

                return GestureDetector(
                  onTap: _showMerchantSelector,
                  child: Container(
                    padding: EdgeInsets.all(48.w),
                    decoration: BoxDecoration(
                      color: AppColors.backgroundLight,
                      borderRadius: BorderRadius.circular(24.r),
                      border: Border.all(color: AppColors.border, width: 1.5),
                    ),
                    child: Row(
                      children: [
                        Icon(
                          Icons.store_rounded,
                          size: 48.sp,
                          color: AppColors.textSecondary,
                        ),
                        SizedBox(width: 32.w),
                        Expanded(
                          child: Text(
                            selectedMerchant?.businessName ??
                                'Select a merchant',
                            style: TextStyle(
                              fontSize: 32.sp,
                              color: selectedMerchant != null
                                  ? AppColors.textPrimary
                                  : AppColors.textHint,
                              fontWeight: selectedMerchant != null
                                  ? FontWeight.w500
                                  : FontWeight.w400,
                            ),
                          ),
                        ),
                        Icon(
                          Icons.arrow_drop_down_rounded,
                          size: 56.sp,
                          color: AppColors.textSecondary,
                        ),
                      ],
                    ),
                  ),
                );
              }),

              SizedBox(height: 48.h),

              Text(
                'Amount',
                style: TextStyle(
                  fontSize: 32.sp,
                  fontWeight: FontWeight.w600,
                  color: AppColors.textPrimary,
                ),
              ),
              SizedBox(height: 24.h),
              CustomTextField(
                controller: _amountController,
                hintText: 'Enter amount',
                keyboardType: TextInputType.numberWithOptions(decimal: true),
                prefixIcon: Icon(
                  Icons.currency_rupee_rounded,
                  size: 48.sp,
                  color: AppColors.textSecondary,
                ),
                validator: _validateAmount,
              ),

              SizedBox(height: 96.h),

              Obx(
                () => CustomButton(
                  text: 'Pay Now',
                  onPressed: _handleCreatePayment,
                  isLoading: _paymentController.isCreatingPayment.value,
                ),
              ),

              SizedBox(height: 48.h),
            ],
          ),
        ),
      ),
    );
  }
}
