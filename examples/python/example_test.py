import unittest
from app import custom_sum

class TestSumFunction(unittest.TestCase):
  def test_sum(self):
    self.assertEqual(custom_sum(1, 1), 2)
    self.assertEqual(custom_sum(10, 10), 20)
    self.assertEqual(custom_sum(100, 100), 200)
    self.assertEqual(custom_sum(1000, 1000), 2000)
