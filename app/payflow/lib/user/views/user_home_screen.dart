import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get.dart';
import 'package:payflow/api/api_repo.dart';
import 'package:payflow/common/controllers/auth_controller.dart';
import 'package:payflow/common/utils/app_colors.dart';
import 'package:payflow/common/views/create_payment_screen.dart';

class UserHomeScreen extends StatefulWidget {
  const UserHomeScreen({super.key});

  @override
  State<UserHomeScreen> createState() => _UserHomeScreenState();
}

class _UserHomeScreenState extends State<UserHomeScreen> {
  final authController = Get.find<AuthController>();
  final apiRepo = ApiRepo();

  final RxBool isLoadingBalance = false.obs;
  final RxBool isLoadingTransactions = false.obs;
  final RxDouble balance = 0.0.obs;
  final RxList transactions = [].obs;

  @override
  void initState() {
    super.initState();
    _loadData();
  }

  Future<void> _loadData() async {
    await Future.wait([_fetchBalance(), _fetchTransactions()]);
  }

  Future<void> _fetchBalance() async {
    try {
      isLoadingBalance.value = true;
      final response = await apiRepo.getWalletBalance();
      balance.value = (response['balance'] as num?)?.toDouble() ?? 0.0;
    } catch (e) {
      print('Fetch balance error: $e');
    } finally {
      isLoadingBalance.value = false;
    }
  }

  Future<void> _fetchTransactions() async {
    try {
      isLoadingTransactions.value = true;
      final response = await apiRepo.getWalletTransactions(limit: 5);
      if (response['transactions'] != null) {
        transactions.value = response['transactions'];
      }
    } catch (e) {
      print('Fetch transactions error: $e');
    } finally {
      isLoadingTransactions.value = false;
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.backgroundLight,
      appBar: AppBar(
        backgroundColor: Colors.white,
        elevation: 0,
        title: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'Home',
              style: TextStyle(
                fontSize: 42.sp,
                fontWeight: FontWeight.bold,
                color: AppColors.textPrimary,
              ),
            ),
            Obx(
              () => Text(
                authController.userEmail.value ?? '',
                style: TextStyle(
                  fontSize: 24.sp,
                  color: AppColors.textSecondary,
                ),
              ),
            ),
          ],
        ),
        actions: [
          IconButton(
            icon: Icon(Icons.logout, size: 48.sp),
            onPressed: () => authController.logout(),
          ),
        ],
      ),
      body: RefreshIndicator(
        onRefresh: _loadData,
        child: SingleChildScrollView(
          physics: const AlwaysScrollableScrollPhysics(),
          padding: EdgeInsets.all(48.w),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              _buildBalanceCard(),
              SizedBox(height: 48.h),
              _buildQuickActions(),
              SizedBox(height: 48.h),
              Text(
                'Recent Transactions',
                style: TextStyle(
                  fontSize: 36.sp,
                  fontWeight: FontWeight.bold,
                  color: AppColors.textPrimary,
                ),
              ),
              SizedBox(height: 24.h),
              _buildTransactionsList(),
            ],
          ),
        ),
      ),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: () async {
          await Get.to(() => const CreatePaymentScreen());
          _loadData();
        },
        backgroundColor: AppColors.primary,
        icon: Icon(Icons.payment, size: 48.sp),
        label: Text(
          'Pay',
          style: TextStyle(fontSize: 32.sp, fontWeight: FontWeight.w600),
        ),
      ),
    );
  }

  Widget _buildBalanceCard() {
    return Container(
      padding: EdgeInsets.all(60.w),
      decoration: BoxDecoration(
        gradient: LinearGradient(
          colors: [AppColors.primary, AppColors.primaryDark],
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
        ),
        borderRadius: BorderRadius.circular(32.r),
        boxShadow: [
          BoxShadow(
            color: AppColors.primary.withOpacity(0.3),
            blurRadius: 20,
            offset: Offset(0, 10),
          ),
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text(
                'Wallet Balance',
                style: TextStyle(
                  fontSize: 32.sp,
                  color: Colors.white.withOpacity(0.9),
                ),
              ),
              Container(
                padding: EdgeInsets.all(24.w),
                decoration: BoxDecoration(
                  color: Colors.white.withOpacity(0.2),
                  borderRadius: BorderRadius.circular(16.r),
                ),
                child: Icon(
                  Icons.account_balance_wallet,
                  size: 48.sp,
                  color: Colors.white,
                ),
              ),
            ],
          ),
          SizedBox(height: 32.h),
          Obx(() {
            if (isLoadingBalance.value) {
              return SizedBox(
                height: 60.h,
                child: Center(
                  child: CircularProgressIndicator(color: Colors.white),
                ),
              );
            }
            return Text(
              '₹${balance.value.toStringAsFixed(2)}',
              style: TextStyle(
                fontSize: 72.sp,
                fontWeight: FontWeight.bold,
                color: Colors.white,
              ),
            );
          }),
        ],
      ),
    );
  }

  Widget _buildQuickActions() {
    return Row(
      children: [
        Expanded(
          child: _buildActionCard(
            icon: Icons.send,
            label: 'Send Money',
            color: AppColors.primary,
            onTap: () async {
              await Get.to(() => const CreatePaymentScreen());
              _loadData();
            },
          ),
        ),
        SizedBox(width: 32.w),
        Expanded(
          child: _buildActionCard(
            icon: Icons.history,
            label: 'History',
            color: AppColors.info,
            onTap: () {},
          ),
        ),
      ],
    );
  }

  Widget _buildActionCard({
    required IconData icon,
    required String label,
    required Color color,
    required VoidCallback onTap,
  }) {
    return InkWell(
      onTap: onTap,
      borderRadius: BorderRadius.circular(24.r),
      child: Container(
        padding: EdgeInsets.all(48.w),
        decoration: BoxDecoration(
          color: Colors.white,
          borderRadius: BorderRadius.circular(24.r),
          border: Border.all(color: AppColors.border),
        ),
        child: Column(
          children: [
            Container(
              padding: EdgeInsets.all(32.w),
              decoration: BoxDecoration(
                color: color.withOpacity(0.1),
                borderRadius: BorderRadius.circular(16.r),
              ),
              child: Icon(icon, size: 56.sp, color: color),
            ),
            SizedBox(height: 24.h),
            Text(
              label,
              style: TextStyle(
                fontSize: 28.sp,
                fontWeight: FontWeight.w600,
                color: AppColors.textPrimary,
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildTransactionsList() {
    return Obx(() {
      if (isLoadingTransactions.value) {
        return Center(
          child: Padding(
            padding: EdgeInsets.all(96.h),
            child: CircularProgressIndicator(),
          ),
        );
      }

      if (transactions.isEmpty) {
        return _buildEmptyTransactions();
      }

      return ListView.separated(
        shrinkWrap: true,
        physics: NeverScrollableScrollPhysics(),
        itemCount: transactions.length,
        separatorBuilder: (_, __) => SizedBox(height: 24.h),
        itemBuilder: (context, index) {
          final txn = transactions[index];
          return _buildTransactionTile(txn);
        },
      );
    });
  }

  Widget _buildTransactionTile(dynamic txn) {
    final isCredit = txn['type'] == 'CREDIT';
    final amount = (txn['amount'] as num?)?.toDouble() ?? 0.0;

    return Container(
      padding: EdgeInsets.all(48.w),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(24.r),
        border: Border.all(color: AppColors.border),
      ),
      child: Row(
        children: [
          Container(
            padding: EdgeInsets.all(24.w),
            decoration: BoxDecoration(
              color: (isCredit ? AppColors.success : AppColors.error)
                  .withOpacity(0.1),
              borderRadius: BorderRadius.circular(12.r),
            ),
            child: Icon(
              isCredit ? Icons.arrow_downward : Icons.arrow_upward,
              size: 40.sp,
              color: isCredit ? AppColors.success : AppColors.error,
            ),
          ),
          SizedBox(width: 32.w),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  isCredit ? 'Received' : 'Sent',
                  style: TextStyle(
                    fontSize: 32.sp,
                    fontWeight: FontWeight.w600,
                    color: AppColors.textPrimary,
                  ),
                ),
                SizedBox(height: 4.h),
                Text(
                  txn['reference'] ?? '',
                  style: TextStyle(
                    fontSize: 24.sp,
                    color: AppColors.textSecondary,
                  ),
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                ),
              ],
            ),
          ),
          Text(
            '${isCredit ? '+' : '-'}₹${amount.toStringAsFixed(2)}',
            style: TextStyle(
              fontSize: 36.sp,
              fontWeight: FontWeight.bold,
              color: isCredit ? AppColors.success : AppColors.error,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildEmptyTransactions() {
    return Container(
      padding: EdgeInsets.all(96.w),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(24.r),
      ),
      child: Column(
        children: [
          Icon(
            Icons.receipt_long_outlined,
            size: 120.sp,
            color: AppColors.textHint,
          ),
          SizedBox(height: 32.h),
          Text(
            'No transactions yet',
            style: TextStyle(fontSize: 32.sp, color: AppColors.textSecondary),
          ),
          SizedBox(height: 16.h),
          Text(
            'Make your first payment to get started',
            style: TextStyle(fontSize: 26.sp, color: AppColors.textHint),
          ),
        ],
      ),
    );
  }
}
