import 'package:flutter_test/flutter_test.dart';

import 'package:vasoolimilaan_app/main.dart';

void main() {
  test('separates true sales, cash-in and net credit change', () {
    final d = DayClose(3000, 2000, 500);
    expect(d.trueSales, closeTo(3500, 1e-9));
    expect(d.cashIn, closeTo(5000, 1e-9));
    expect(d.netCreditChange, closeTo(-1500, 1e-9));
  });

  testWidgets('renders day-close', (tester) async {
    await tester.pumpWidget(const VasoolimilaanApp());
    expect(find.text('True sales today'), findsOneWidget);
  });
}
