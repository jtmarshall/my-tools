from sqlalchemy import Column, DateTime, String, Integer, ForeignKey, func, update, create_engine
from sqlalchemy.orm import relationship, backref, sessionmaker
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.types import *


Base = declarative_base()

# Connection Setup
host = "your_host"
db_name = "your_db_name"
user = "your_user"
password = "your_pass"
connect_string = 'mysql+pymysql://%s:%s@%s/%s' % (user, password, host, db_name)
engine = create_engine(connect_string)


# Create New Crawl Table
class Crawl(Base):
    __tablename__ = 'crawl'
    id = Column(Integer, primary_key=True)
    start_time = Column(DateTime)
    end_time = Column(DateTime)
    total_crawled = Column(Integer)


# Create New Page Table
class Page(Base):
    __tablename__ = 'page'
    id = Column(Integer, primary_key=True)
    domain = Column(String(32))
    url = Column(String(191))
    # Use default=func.now() to set the default time of a Page to be the current time when the page record was created
    datetime = Column(DateTime, default=func.now())
    status_code = Column(Integer)
    response_time = Column(Float)
    ttfb = Column(Float)
    error = Column(String(32))
    redirects = Column(Integer)
    crawl_id = Column(Integer, ForeignKey('crawl.id'))
    # Use cascade='delete,all' to propagate the deletion of a Department onto its Employees
    crawl = relationship(Crawl, backref=backref('crawl', uselist=True, cascade='delete,all'))


# Load Status table
class Status(Base):
    __tablename__ = 'status'
    __table_args__ = {
        'autoload': True,
        'autoload_with': engine
    }


# Load Outage table
class Outage(Base):
    __tablename__ = 'outages'
    __table_args__ = {
        'autoload': True,
        'autoload_with': engine
    }


def loadSession():
    metadata = Base.metadata
    Session = sessionmaker(bind=engine)
    session = Session()
    return session
