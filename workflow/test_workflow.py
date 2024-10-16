# coding=utf-8
import unittest

import korean_search

class MyTestCase(unittest.TestCase):
    def test_naver_search(self):
        res = korean_search.get_data('사랑')

        self.assertTrue(len(res['items']) > 0)

if __name__ == '__main__':
    unittest.main()
