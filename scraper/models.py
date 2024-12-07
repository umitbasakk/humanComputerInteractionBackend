

class RequestTweet:
    def __init__(self,startDate,endDate,tag):
        self.startDate=startDate
        self.endDate = endDate
        self.tag = tag
    
    def to_dict(self):
        return {
            "startDate": self.start,
            "endDate": self.end,
            "tag": self.tag
        }
    
class Tweet:
    def __init__(self,publishDate,tweet,classify):
        self.publishDate=publishDate
        self.tweet = tweet
        self.classify = classify
    
    def to_dict(self):
        return {
            "publishDate": self.publishDate,
            "tweet": self.tweet,
            "classify": self.classify
        }