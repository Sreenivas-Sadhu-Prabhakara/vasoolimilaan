import 'package:flutter/material.dart';

void main() => runApp(const VasoolimilaanApp());

/// Vasoolimilaan — day-close that separates today's cash sales from dues
/// collected and new credit given, so recovery/fresh-credit don't distort
/// "did I make money today". Mirrors the Go journal service.
class VasoolimilaanApp extends StatelessWidget {
  const VasoolimilaanApp({super.key});
  @override
  Widget build(BuildContext context) => MaterialApp(
        title: 'Vasoolimilaan',
        debugShowCheckedModeBanner: false,
        theme: ThemeData(colorSchemeSeed: const Color(0xFF3E7E6E), useMaterial3: true),
        home: const HomePage(),
      );
}

class DayClose {
  final double cashSales, duesCollected, newCredit;
  DayClose(this.cashSales, this.duesCollected, this.newCredit);
  double get trueSales => cashSales + newCredit;
  double get cashIn => cashSales + duesCollected;
  double get netCreditChange => newCredit - duesCollected;
}

class HomePage extends StatefulWidget {
  const HomePage({super.key});
  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  final _cash = TextEditingController(text: '3000');
  final _dues = TextEditingController(text: '2000');
  final _credit = TextEditingController(text: '500');

  double _n(TextEditingController c) => double.tryParse(c.text.trim()) ?? 0;

  @override
  Widget build(BuildContext context) {
    final d = DayClose(_n(_cash), _n(_dues), _n(_credit));
    String m(double v) => '₹${v.toStringAsFixed(2)}';
    return Scaffold(
      appBar: AppBar(
        title: const Text('Vasoolimilaan · day-close'),
        backgroundColor: Theme.of(context).colorScheme.primaryContainer,
      ),
      body: ListView(padding: const EdgeInsets.all(16), children: [
        _f(_cash, 'Cash sales today ₹'),
        _f(_dues, 'Dues collected (old) ₹'),
        _f(_credit, 'New credit given ₹'),
        const SizedBox(height: 16),
        Card(
          color: Theme.of(context).colorScheme.primaryContainer,
          child: Padding(padding: const EdgeInsets.all(16),
            child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
              _big('True sales today', m(d.trueSales)),
              const Divider(),
              _row('Cash into till', m(d.cashIn)),
              _row('Net credit change', '${d.netCreditChange >= 0 ? '+' : ''}${m(d.netCreditChange)}'),
              const SizedBox(height: 8),
              const Text('A till fattened by dues recovery is not a great sales day; thinned by fresh credit is not a loss.',
                  style: TextStyle(fontSize: 12)),
            ])),
        ),
      ]),
    );
  }

  Widget _big(String k, String v) => Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
        Text(k), Text(v, style: const TextStyle(fontSize: 28, fontWeight: FontWeight.bold)),
      ]);

  Widget _row(String k, String v) => Padding(
        padding: const EdgeInsets.symmetric(vertical: 3),
        child: Row(mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [Text(k), Text(v, style: const TextStyle(fontWeight: FontWeight.w600))]),
      );

  Widget _f(TextEditingController c, String label) => Padding(
        padding: const EdgeInsets.symmetric(vertical: 6),
        child: TextField(controller: c, keyboardType: const TextInputType.numberWithOptions(decimal: true),
          decoration: InputDecoration(labelText: label, border: const OutlineInputBorder()),
          onChanged: (_) => setState(() {})),
      );
}
