import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:slides/fetch.dart';

class SliderPage extends StatefulWidget {
  final String code;

  const SliderPage({required this.code, super.key});

  @override
  State<SliderPage> createState() => _SliderPageState();
}

class _SliderPageState extends State<SliderPage> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Slides Controller')),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            CupertinoButton.filled(
              onPressed: () {
                var res = dio.post("/api/next", data: {'code': widget.code});
                res.then((value) => debugPrint("onNext: ${value.statusCode} ${value.data}"));
              },
              child: const Text("Prev"),
            ),
            CupertinoButton.filled(
              onPressed: () {
                var res = dio.post("/api/prev", data: {'code': widget.code});
                res.then((value) => debugPrint("onNext: ${value.statusCode} ${value.data}"));
              },
              child: const Text("Next"),
            )
          ]
        )
      ),
    );
  }
}